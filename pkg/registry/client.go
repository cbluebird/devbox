package registry

import (
	"bytes"
	"errors"
	"github.com/google/go-containerregistry/pkg/name"
	"io"
	"log"
	"net/http"
	"time"

	retry "github.com/avast/retry-go"
)

type Client struct {
	Username string
	Password string
}

var (
	ErrorManifestNotFound = errors.New("manifest not found")
)

type ImageOptions struct {
	HostName  string
	ImageName string
	Tag       string
}

func (c *Client) Tag(originOpt, newOpt ImageOptions) error {
	return retry.Do(func() error {
		manifest, err := c.pullManifest(originOpt)
		if err != nil {
			log.Println("failed to pull manifest", err)
			return err
		}
		return c.pushManifest(newOpt, manifest)
	}, retry.Delay(time.Second*5), retry.Attempts(3), retry.LastErrorOnly(true))
}

func (c *Client) pullManifest(imageOpt ImageOptions) ([]byte, error) {
	client := http.DefaultClient
	url := "https://" + imageOpt.HostName + "/v2/" + imageOpt.ImageName + "/manifests/" + imageOpt.Tag
	log.Println("pulling manifest", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("failed to pull manifest", err)
		return nil, err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrorManifestNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to pull manifest: " + resp.Status)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyText, nil
}

func (c *Client) pushManifest(options ImageOptions, manifest []byte) error {
	client := http.DefaultClient
	url := "https://" + options.HostName + "/v2/" + options.ImageName + "/manifests/" + options.Tag

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(manifest))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("failed to push manifest", err)
		return err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusCreated {
		return errors.New("failed to push manifest: " + resp.Status)
	}

	return nil
}

func GetImageInfo(imageRef string) (*ImageOptions, error) {
	res, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, err
	}
	repo := res.Context()
	log.Println(repo.RegistryStr())
	log.Println(repo.RepositoryStr())
	log.Println(res.Identifier())

	return &ImageOptions{
		HostName:  repo.RegistryStr(),
		ImageName: repo.RepositoryStr(),
		Tag:       res.Identifier(),
	}, nil
}

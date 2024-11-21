package registry

import (
	"log"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/heroku/docker-registry-client/registry"
)

func Tag(original, target string) error {
	return nil
}

type Client struct {
	Username string
	Password string
}

type ImageOptions struct {
	HostName  string
	ImageName string
	Tag       string
}

func (c *Client) ReTag(originImage, newImage string) error {

	originOpt, err := GetImageInfo(originImage)
	if err != nil {
		log.Fatalf("Error getting origin image info: %v", err)
	}
	newOpt, err := GetImageInfo(newImage)
	if err != nil {
		log.Fatalf("Error getting new image info: %v", err)
	}

	hub, err := registry.New("http://"+originOpt.HostName, c.Username, c.Password)
	if nil != err {
		log.Println("failed to create hub", err)
		return err
	}
	manifest, err := hub.ManifestV2(originOpt.ImageName, originOpt.Tag)
	if nil != err {
		log.Println("failed to get manifest", err)
		return err
	}
	err = hub.PutManifest(newOpt.ImageName, newOpt.Tag, manifest)
	if err != nil {
		log.Println("failed to put manifest", err)
		return err
	}
	log.Println("Successfully pushed manifest to", newOpt.ImageName)
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

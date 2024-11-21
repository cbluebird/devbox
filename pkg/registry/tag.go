package tag

import (
	"log/slog"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/heroku/docker-registry-client/registry"
)

var TagClient *Client

type Client struct {
	Username string
	Password string
}

type ImageOptions struct {
	HostName  string
	ImageName string
	Tag       string
}

func Init(user, password string) {
	TagClient = &Client{
		Username: user,
		Password: password,
	}
}

func (c *Client) Tag(originImage, newImage string) error {
	originOpt, err := GetImageInfo(originImage)
	if err != nil {
		slog.Error("Error getting origin image info", "Error", err)
		return err
	}
	newOpt, err := GetImageInfo(newImage)
	if err != nil {
		slog.Error("Error getting new image info", "Error", err)
		return err
	}

	hub, err := registry.New("http://"+originOpt.HostName, c.Username, c.Password)
	if nil != err {
		slog.Error("Failed to create hub", "Error", err)
		return err
	}
	manifest, err := hub.ManifestV2(originOpt.ImageName, originOpt.Tag)
	if nil != err {
		slog.Error("Failed to get manifest", "Error", err)
		return err
	}
	err = hub.PutManifest(newOpt.ImageName, newOpt.Tag, manifest)
	if err != nil {
		slog.Error("Failed to put manifest", "Error", err)
		return err
	}
	slog.Info("Tag success", "OriginImage", originImage, "NewImage", newImage)
	return nil
}

func GetImageInfo(imageRef string) (*ImageOptions, error) {
	res, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, err
	}
	repo := res.Context()
	return &ImageOptions{
		HostName:  repo.RegistryStr(),
		ImageName: repo.RepositoryStr(),
		Tag:       res.Identifier(),
	}, nil
}

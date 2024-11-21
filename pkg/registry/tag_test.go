package registry

import (
	"log"
	"testing"
)

func TestTag(t *testing.T) {
	originOpt, err := GetImageInfo("docker.io/cbluebird/crk-nginx:latest")
	if err != nil {
		log.Println(err)
	}

	client := Client{
		Username: "",
		Password: "",
	}

	newOpt := ImageOptions{
		HostName:  "index.docker.io",
		ImageName: "cbluebird/crk-nginx-test",
		Tag:       "test",
	}

	if err := client.Tag(*originOpt, newOpt); err != nil {
		log.Println(err)
	}
}

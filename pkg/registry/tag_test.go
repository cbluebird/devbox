package registry

import (
	"testing"
)

func TestTag(t *testing.T) {
	client := &Client{
		Username: "",
		Password: "",
	}
	err := client.ReTag("sealos.hub:5000/cbluebird/crk-nginx-test:latest-1", "sealos.hub:5000/cbluebird/crk-nginx-test:latest-2")
	if err != nil {
		return
	}
}

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labring/sealos/service/devbox/pkg/registry"
)

type TagRequest struct {
	Original string `json:"original"`
	Target   string `json:"target"`
}

type TagResponse struct {
}

func Tag(c *gin.Context) {
	var request TagRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := registry.Tag(request.Original, request.Target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, TagResponse{})
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labring/sealos/service/devbox/api"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	r := gin.Default()
	r.POST("/tag", api.Tag)
	r.Run(":8080")
}

package main

import (
	"flag"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/labring/sealos/service/devbox/api"
	tag "github.com/labring/sealos/service/devbox/pkg/registry"
)

func main() {
	user := flag.String("user", "admin", "Username for authentication")
	password := flag.String("password", "password", "Password for authentication")
	DebugFlag := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	tag.Init(*user, *password)
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if *DebugFlag {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	r := gin.Default()
	r.POST("/tag", api.Tag)
	r.Run(":8080")
}

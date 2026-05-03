package main

import (
	"go-fake-flix/internal/apierrors"
	"go-fake-flix/internal/modules/video"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.Use(apierrors.ErrorHandler())
	public := server.Group("/api/v1")
	video.RegisterVideoRoutes(public)

	server.Run(":8080")
}

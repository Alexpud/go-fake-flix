package main

import (
	"net/http"
	"path/filepath"

	"go-fake-flix/internal/apierrors"

	"github.com/gin-gonic/gin"
)

var db = map[string]string{
	"hello": "Hello, World!",
}

func main() {
	server := gin.Default()
	server.Use(apierrors.ErrorHandler())
	public := server.Group("/api/v1")

	videoController := public.Group("/video")

	videoController.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.Error(&apierrors.AppError{Code: "FILE_NOT_FOUND", Message: "File not found"})
			return
		}

		dst := filepath.Join("./files/", filepath.Base(file.Filename))
		c.SaveUploadedFile(file, dst)

		db[file.Filename] = dst

		c.Status(http.StatusOK)
	})

	videoController.GET("/:id", func(c *gin.Context) {
		// There should be logic to handle different ranges in the video: getting the video from a starting point and stuff
		video, ok := db[c.Param("id")]
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.File(video)
	})

	server.Run(":8080")
}

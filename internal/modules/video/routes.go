package video

import (
	"go-fake-flix/internal/apierrors"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var db = map[string]string{
	"hello": "Hello, World!",
}

func RegisterVideoRoutes(router *gin.RouterGroup) {
	videoRouter := router.Group("/video")
	videoRouter.POST("/upload", uploadVideo)
	videoRouter.GET("/:id", getVideo)
}

func uploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(&apierrors.AppError{Code: "FILE_NOT_FOUND", Message: "File not found"})
		return
	}

	dst := filepath.Join("./files/", filepath.Base(file.Filename))
	c.SaveUploadedFile(file, dst)
	db[file.Filename] = dst
	c.Status(http.StatusOK)
}

func getVideo(c *gin.Context) {
	// There should be logic to handle different ranges in the video: getting the video from a starting point and stuff
	video, ok := db[c.Param("id")]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	c.File(video)
}

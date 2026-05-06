package video

import (
	"go-fake-flix/internal/apierrors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var db = map[string]string{
	"hello": "Hello, World!",
}

func RegisterVideoRoutes(r *chi.Mux) {
	r.Route("/api/v1/content", func(r chi.Router) {
		r.Post("/upload", uploadVideo)
		r.Get("/stream/{id}", getVideo)
	})
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &apierrors.AppError{Code: "FILE_NOT_FOUND", Message: "File not found"})
		return
	}

	dst := filepath.Join("./files/", filepath.Base(r.FormValue("filename")))
	dstFile, err := os.Create(dst)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, &apierrors.AppError{Code: "FILE_NOT_CREATED", Message: "File not created"})
		return
	}
	defer dstFile.Close()
	io.Copy(dstFile, file)
	db[r.FormValue("filename")] = dst
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &apierrors.AppError{Code: "VIDEO_UPLOADED", Message: "Video uploaded"})
}

func getVideo(w http.ResponseWriter, r *http.Request) {
	// There should be logic to handle different ranges in the video: getting the video from a starting point and stuff
	video, ok := db[r.URL.Query().Get("id")]
	if !ok {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, &apierrors.AppError{Code: "VIDEO_NOT_FOUND", Message: "Video not found"})
		return
	}
	http.ServeFile(w, r, video)
}

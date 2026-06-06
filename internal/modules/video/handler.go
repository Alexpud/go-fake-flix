package video

import (
	"go-fake-flix/internal/apierrors"
	"go-fake-flix/internal/modules/video/usecases"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	MaxFileSize = 5 * 1024 * 1024 // 5MB
)

// allowedUploadExtensions are suffixes including the dot, lowercased (see parseMultipartUpload).
var allowedUploadExtensions = []string{".mp4", ".webm", ".mkv"}

type ResponseSuccess struct {
	Message string `json:"message,omitempty"`
}

var db = map[string]string{
	"hello": "Hello, World!",
}

func RegisterVideoRoutes(r *chi.Mux) {
	r.Route("/api/v1/content", func(r chi.Router) {
		r.Post("/upload", uploadVideo)
		r.Get("/stream/{id}", getVideo)
		r.Get("/doc", getDoc)
	})
}

type RouteInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type ResponseDoc struct {
	Routes []RouteInfo `json:"routes"`
}

func getDoc(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, ResponseDoc{Routes: []RouteInfo{
		{Method: "POST", Path: "/api/v1/content/upload", Description: "Upload a video file."},
		{Method: "GET", Path: "/api/v1/content/stream/{id}", Description: "Stream a video by id."},
		{Method: "GET", Path: "/api/v1/content/doc", Description: "Get API documentation for video routes."},
	}})
}

// @
// @Summary      Upload a video file
// @Description  Upload a video file
// @Tags         Video
// @Accept       multipart/form-data
// @Produce      json
// @Param        file    formData file    true  "Video file"
// @Param        filename formData string true  "File name"
// @Success      200
// @Router       /api/v1/content/upload [post]
func uploadVideo(w http.ResponseWriter, r *http.Request) {
	file, fileName, err := parseMultipartUpload(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err)
		return
	}
	defer file.Close()

	filePath, errUsecase := usecases.UploadVideo(file, fileName)
	if errUsecase != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, errUsecase)
		return
	}

	db[fileName] = filePath
	render.Status(r, http.StatusOK)
}

func parseMultipartUpload(r *http.Request) (multipart.File, string, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, "", &apierrors.AppError{Code: "FILE_NOT_FOUND", Message: "File not found"}
	}

	if fileHeader.Size > MaxFileSize {
		return nil, "", &apierrors.AppError{Code: "FILE_TOO_LARGE", Message: "Max file size is 5MB"}
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" || !slices.Contains(allowedUploadExtensions, ext) {
		return nil, "", &apierrors.AppError{
			Code:    "FILE_EXTENSION_INVALID",
			Message: "Invalid file extension. Allowed: " + strings.Join(allowedUploadExtensions, ", "),
		}
	}

	return file, fileHeader.Filename, nil
}

// Exemple Doc
// @Summary      Get stream content
// @Description  Get stream content
// @Tags         Video
// @Produce      json
// @Success      200				 {object}  ResponseSuccess
// @Router       /api/v1/content/stream/{id} [get]
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

package video

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"go-fake-flix/internal/apierrors"
	"go-fake-flix/internal/modules/video/repository"
	"go-fake-flix/internal/modules/video/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	MaxFileSize = 5 * 1024 * 1024
)

var allowedUploadExtensions = []string{".mp4", ".webm", ".mkv"}

type ResponseSuccess struct {
	Message string `json:"message,omitempty"`
}

func RegisterVideoRoutes(r *chi.Mux, repo repository.VideoRepository) {
	r.Route("/api/v1/content", func(r chi.Router) {
		r.Post("/upload", uploadVideo(repo))
		r.Get("/stream/{id}", getVideo(repo))
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

// @Summary      Upload a video file
// @Description  Upload a video file
// @Tags         Video
// @Accept       multipart/form-data
// @Produce      json
// @Param        file    formData file    true  "Video file"
// @Param        filename formData string true  "File name"
// @Success      200
// @Router       /api/v1/content/upload [post]
func uploadVideo(repo repository.VideoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, fileName, appErr := parseMultipartUpload(r)
		if appErr != nil {
			renderAppError(w, r, appErr.Status, appErr)
			return
		}
		defer file.Close()

		if appErr := usecases.UploadVideo(r.Context(), repo, fileName, fileName); appErr != nil {
			renderAppError(w, r, 400, appErr)
			return
		}

		render.Status(r, http.StatusOK)
	}
}

func parseMultipartUpload(r *http.Request) (multipart.File, string, *apierrors.ProblemDetails) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, "", apierrors.CreateProblemDetails("FILE_NOT_FOUND", "File not found", http.StatusBadRequest)
	}

	if fileHeader.Size > MaxFileSize {
		return nil, "", apierrors.CreateProblemDetails("FILE_TOO_LARGE", "Max file size is 5MB", http.StatusBadRequest)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" || !slices.Contains(allowedUploadExtensions, ext) {
		return nil, "", apierrors.CreateProblemDetails(
			"FILE_EXTENSION_INVALID",
			"Invalid file extension. Allowed: "+strings.Join(allowedUploadExtensions, ", "),
			http.StatusBadRequest,
		)
	}

	return file, fileHeader.Filename, nil
}

// @Summary      Get stream content
// @Description  Get stream content
// @Tags         Video
// @Produce      json
// @Success      200				 {object}  ResponseSuccess
// @Router       /api/v1/content/stream/{id} [get]
func getVideo(repo repository.VideoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		v, err := usecases.GetVideo(r.Context(), repo, id)
		if err != nil {
			renderAppError(w, r, 400, err)
			return
		}

		http.ServeFile(w, r, v.FilePath)
	}
}

func renderAppError(w http.ResponseWriter, r *http.Request, status int, appErr any) {
	render.Status(r, status)
	render.JSON(w, r, appErr)
}

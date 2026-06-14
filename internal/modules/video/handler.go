package video

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"go-fake-flix/internal/apierrors"
	"go-fake-flix/internal/modules/video/service"

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

func RegisterVideoRoutes(r *chi.Mux, service service.VideoService, mediaService *service.MediaService) {
	r.Route("/api/v1/content", func(r chi.Router) {
		r.Post("/upload", uploadVideo(service))
		r.Get("/stream/{id}", getVideo(service))
		r.Get("/media/{id}", streamMedia(mediaService))
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
		{Method: "GET", Path: "/api/v1/content/media/{id}", Description: "Stream a video byte range by id."},
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
func uploadVideo(service service.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, fileName, appErr := parseMultipartUpload(r)
		if appErr != nil {
			renderAppError(w, r, appErr.Status, appErr)
			return
		}
		defer func() { _ = file.Close() }()

		videoId, busErr := service.UploadVideo(r.Context(), fileName, file)
		if busErr != nil {
			renderAppError(w, r, 400, busErr)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, videoId)
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
func getVideo(service service.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		i, err :=strconv.ParseInt(id, 10, 64)
		if err != nil {
			renderAppError(w,r, 400, apierrors.CreateProblemDetails(
				"INVALID_VIDEO_ID",
				"VideoId must be an int64",
				http.StatusBadRequest,
			))
			return
		}
		
		v, err := service.GetVideo(r.Context(), i)
		if err != nil {
			renderAppError(w, r, 400, err)
			return
		}

		http.ServeFile(w, r, v.FilePath)
	}
}

// @Summary      Stream video content
// @Description  Stream a video file. Supports HTTP byte range requests (Range header) for partial content.
// @Tags         Video
// @Produce      application/octet-stream
// @Param        id    path     int     true  "Video ID"
// @Param        Range header   string false "Byte range (e.g. bytes=0-1023)"
// @Success      200
// @Success      206
// @Router       /api/v1/content/media/{id} [get]
func streamMedia(mediaService *service.MediaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			renderAppError(w, r, http.StatusBadRequest, apierrors.CreateProblemDetails(
				"INVALID_VIDEO_ID", "Video ID must be an integer", http.StatusBadRequest,
			))
			return
		}

		rangeStr := r.Header.Get("Range")
		if rangeStr == "" {
			data, busErr := mediaService.GetMedia(r.Context(), id, nil)
			if busErr != nil {
				renderAppError(w, r, http.StatusBadRequest, busErr)
				return
			}
			w.Header().Set("Content-Type", "video/mp4")
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
			return
		}

		vr, appErr := parseRange(rangeStr)
		if appErr != nil {
			renderAppError(w, r, http.StatusBadRequest, appErr)
			return
		}

		data, busErr := mediaService.GetMedia(r.Context(), id, vr)
		if busErr != nil {
			renderAppError(w, r, http.StatusBadRequest, busErr)
			return
		}

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/*", vr.Start, vr.End))
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(data)
	}
}

func parseRange(rangeStr string) (*service.VideoRange, *apierrors.ProblemDetails) {
	const prefix = "bytes="
	if !strings.HasPrefix(rangeStr, prefix) {
		return nil, apierrors.CreateProblemDetails(
			"INVALID_RANGE_FORMAT", "Range must start with 'bytes='", http.StatusBadRequest,
		)
	}

	parts := strings.SplitN(strings.TrimPrefix(rangeStr, prefix), "-", 2)
	if len(parts) != 2 {
		return nil, apierrors.CreateProblemDetails(
			"INVALID_RANGE_FORMAT", "Range must be in format bytes=start-end", http.StatusBadRequest,
		)
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, apierrors.CreateProblemDetails(
			"INVALID_RANGE_START", "Range start must be an integer", http.StatusBadRequest,
		)
	}

	end, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, apierrors.CreateProblemDetails(
			"INVALID_RANGE_END", "Range end must be an integer", http.StatusBadRequest,
		)
	}

	if start > end {
		return nil, apierrors.CreateProblemDetails(
			"INVALID_RANGE", "Range start must be less than or equal to end", http.StatusBadRequest,
		)
	}

	return &service.VideoRange{Start: start, End: end}, nil
}

func renderAppError(w http.ResponseWriter, r *http.Request, status int, appErr any) {
	render.Status(r, status)
	render.JSON(w, r, appErr)
}

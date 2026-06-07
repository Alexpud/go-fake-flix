package main

import (
	"fmt"
	"go-fake-flix/internal/modules/video"
	videoRepository "go-fake-flix/internal/modules/video/repository"
	"go-fake-flix/internal/openapi"
	"log"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title           Fake flix API
// @version         1.0
// @description     Fake flix API
// @termsOfService  http://swagger.io/terms/

// @BasePath  /
func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	videoRepo := videoRepository.New()
	video.RegisterVideoRoutes(r, videoRepo)

	r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	r.Get("/openapi.json", openapi.Handler("docs/swagger.json"))

	r.Get("/reference", func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		} else if xfProto := r.Header.Get("X-Forwarded-Proto"); xfProto != "" {
			scheme = xfProto
		}

		// Use an absolute URL so Scalar doesn't try to treat "/openapi.json" as a Windows file path.
		specURL := fmt.Sprintf("%s://%s/openapi.json", scheme, r.Host)

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: specURL,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: false,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, htmlContent)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

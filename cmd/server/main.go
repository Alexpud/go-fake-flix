package main

import (
	"go-fake-flix/internal/modules/video"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var router = chi.NewRouter()

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// server := gin.Default()
	// server.Use(apierrors.ErrorHandler())
	// public := server.Group("/api/v1")
	video.RegisterVideoRoutes(r)

	http.ListenAndServe(":8080", r)
}

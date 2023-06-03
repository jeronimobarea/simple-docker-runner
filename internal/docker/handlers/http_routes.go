package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/jeronimobarea/simple-docker-runner/internal/docker"
)

func RegisterRoutes(r chi.Router, dockerSvc docker.Service) {
	h := NewHandler(dockerSvc)

	r.Post("/v1/docker/run", h.Run)
}

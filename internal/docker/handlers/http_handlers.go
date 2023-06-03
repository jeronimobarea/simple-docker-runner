package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jeronimobarea/simple-docker-runner/internal/docker"
)

type handler struct {
	dockerSvc docker.Service
}

func NewHandler(dockerSvc docker.Service) *handler {
	return &handler{dockerSvc}
}

func (h *handler) Run(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var cmd docker.RunCMD
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.dockerSvc.Run(ctx, cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

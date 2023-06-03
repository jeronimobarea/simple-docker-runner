package handlers

import (
	"net/http"
)

type handler struct{}

func (h *handler) Htck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Healthy!!!`))
}

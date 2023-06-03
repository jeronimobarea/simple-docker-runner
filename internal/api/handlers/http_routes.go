package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	h := handler{}

	r.Get("/htck", h.Htck)
}

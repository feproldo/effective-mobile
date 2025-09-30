package subscriptions

import (
	"encoding/json"
	"net/http"

	"github.com/feproldo/effective-mobile/internal/dto"
	subsService "github.com/feproldo/effective-mobile/internal/services/subscriptions"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	services *subsService.Services
}

func NewHandler(services *subsService.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.services.List(r.Context())
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.Subscription
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = h.services.Create(r.Context(), body)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("created"))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam("id")
	list, err := h.services.Get(r.Context())
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

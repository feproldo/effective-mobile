package subscriptions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/feproldo/effective-mobile/internal/dto"
	subsService "github.com/feproldo/effective-mobile/internal/services/subscriptions"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

// @Summary      Get list of the subscriptions
// @Description  Get list of the subscriptions
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   dto.Subscription
// @Failure      404  string    "Subscriptions list is empty"
// @Failure      500  string		"Internal error"
// @Router       /subscriptions [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.services.List(r.Context())
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if len(*list) == 0 {
		log.Info().Msg("subscriptions list is empty")
		http.Error(w, "no content", http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// @Summary      Add a new subscription
// @Description  Add a new subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      plain
// @Param        request body      dto.Subscription true "Subscription data"
// @Success      201    string     "created"
// @Failure      400    string     "bad request"
// @Failure      500    string     "interbal server error"
// @Router       /subscriptions    [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.Subscription
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Error().Err(err).Msg("can't decode request body")
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
	w.WriteHeader(201)
	w.Write([]byte("created"))
}

// @Summary      Get subscription by id
// @Description  Get subscription by it's serial primary key
// @Tags         subscriptions
// @Produce      json
// @Param        id      path        int true "Serial primary key"
// @Success      200     {object}    dto.Subscription
// @Failure      400     string     "Id not found"
// @Failure      404     string     "Subscription not found"
// @Failure      500     string		  "Internal error"
// @Router       /subscriptions/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idParsed, err := strconv.Atoi(id)

	if err != nil {
		log.Err(err).Msg("Can't get URL param \"id\"")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	sub, err := h.services.Get(r.Context(), int32(idParsed))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("subscription not found")
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// @Summary      Get subscription by user_id
// @Description  Get subscription by user_id
// @Tags         subscriptions
// @Produce      json
// @Param        user_id path        string true "user UUID"
// @Success      200     {object}    dto.Subscription
// @Failure      400     string     "user_id (UUID) not found"
// @Failure      404     string     "Subscription not found"
// @Failure      500     string		  "Internal error"
// @Router       /subscriptions/user/{user_id} [get]
func (h *Handler) GetByUserId(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		log.Err(err).Msg("Can't get URL param \"user_id\"")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	list, err := h.services.GetByUserId(r.Context(), userUUID)

	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(*list) == 0 {
		log.Info().Msg("subscriptions list is empty")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// @Summary      Delete subscription by its id
// @Description  Delete subscription by its Serial Primary Key
// @Tags         subscriptions
// @Produce      plain
// @Param        id      path        int true "Serial primary key"
// @Success      204
// @Failure      400     string     "user_id (UUID) not found"
// @Failure      500     string		  "Internal error"
// @Router       /subscriptions/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idParsed, err := strconv.Atoi(id)

	if err != nil {
		log.Err(err).Msg("can't get URL param \"id\"")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.services.Delete(r.Context(), int32(idParsed))

	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Update subscription by its id
// @Description  Update subscription by its Serial Primary Key
// @Tags         subscriptions
// @Produce      plain
// @Param        id      path        int true "Serial primary key"
// @Param        request body        dto.Subscription true "Subscription data"
// @Success      204
// @Failure      400     string     "user_id (UUID) not found"
// @Failure      500     string		  "Internal error"
// @Router       /subscriptions/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idParsed, err := strconv.Atoi(id)

	if err != nil {
		log.Err(err).Msg("can't get URL param \"id\"")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var body dto.Subscription
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Err(err).Msg("can't decode request body")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.services.Update(r.Context(), int32(idParsed), body)

	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Get subscription by id
// @Description  Get subscription by it's serial primary key
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query       string false "user_id (UUID)"
// @Param        servuce_name query       string false "Service name"
// @Param        start_date   query       string false "Start date (MM-YYYY)"
// @Param        end_date     query       string false "End date (MM-YYYY)"
// @Success      200     plain       "Sum"
// @Failure      500     string		  "Internal error"
// @Router       /subscriptions/sum [get]
func (h *Handler) Sum(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	sum, err := h.services.Sum(r.Context(), startDate, endDate, userId, serviceName)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(strconv.Itoa(*sum)))
}

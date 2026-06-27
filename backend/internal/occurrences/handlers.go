package occurrences

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/go-sinmos/internal/auth"
	"github.com/lucaserm/go-sinmos/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(router *chi.Mux, authMiddleware *auth.Middleware) {
	occurrencesRouter := chi.NewRouter()
	occurrencesRouter.Get("/", h.getOccurrences)
	occurrencesRouter.Get("/{id}", h.getOccurrenceByID)
	occurrencesRouter.Put("/{id}", h.updateOccurrence)
	occurrencesRouter.Delete("/{id}", h.deleteOccurrence)

	occurrencesRouter.Post("/", authMiddleware.RequiresAuth(h.createOccurrence))

	router.Mount("/occurrences", occurrencesRouter)
}

func (h *handler) createOccurrence(w http.ResponseWriter, r *http.Request, user repo.User) {
	var payload CreateOccurrencePayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.createOccurrence(r.Context(), payload, user.ID.Bytes)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"occurrence": response,
	})
}

func (h *handler) getOccurrences(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	// Convert offset and limit to int32
	offsetInt := json.ParseInt32(offset, 0)
	limitInt := json.ParseInt32(limit, 10)

	response, err := h.service.getOccurrences(r.Context(), offsetInt, limitInt)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"occurrences": response,
	})
}

func (h *handler) getOccurrenceByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	response, err := h.service.getOccurrenceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrOccurrenceNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"occurrence": response,
	})
}

func (h *handler) updateOccurrence(w http.ResponseWriter, r *http.Request) {
	var payload UpdateOccurrencePayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	response, err := h.service.updateOccurrence(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrOccurrenceNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"occurrence": response,
	})
}

func (h *handler) deleteOccurrence(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	err := h.service.deleteOccurrence(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrOccurrenceNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "occurrence deleted successfully",
	})
}

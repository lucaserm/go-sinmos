package occurrencestypes

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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

func (h *handler) RegisterRoutes(router *chi.Mux) {
	occurrenceTypesRouter := chi.NewRouter()
	occurrenceTypesRouter.Post("/", h.createOccurrenceType)
	occurrenceTypesRouter.Get("/", h.getOccurrenceTypes)
	occurrenceTypesRouter.Get("/{id}", h.getOccurrenceTypeByID)
	occurrenceTypesRouter.Put("/{id}", h.updateOccurrenceType)
	occurrenceTypesRouter.Delete("/{id}", h.deleteOccurrenceType)

	router.Mount("/occurrence-types", occurrenceTypesRouter)
}

func (h *handler) createOccurrenceType(w http.ResponseWriter, r *http.Request) {
	var payload CreateOccurrenceTypePayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.createOccurrenceType(r.Context(), payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"occurrenceType": response,
	})
}

func (h *handler) getOccurrenceTypes(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	// Convert offset and limit to int32
	offsetInt := json.ParseInt32(offset, 0)
	limitInt := json.ParseInt32(limit, 10)

	response, err := h.service.getOccurrenceTypes(r.Context(), offsetInt, limitInt)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"occurrenceTypes": response,
	})
}

func (h *handler) getOccurrenceTypeByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	response, err := h.service.getOccurrenceTypeByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrOccurrenceTypeNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"occurrenceType": response,
	})
}

func (h *handler) updateOccurrenceType(w http.ResponseWriter, r *http.Request) {
	var payload UpdateOccurrenceTypePayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id := chi.URLParam(r, "id")
	response, err := h.service.updateOccurrenceType(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrOccurrenceTypeNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"occurrenceType": response,
	})
}

func (h *handler) deleteOccurrenceType(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.deleteOccurrenceType(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrOccurrenceTypeNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "occurrence type deleted successfully",
	})
}

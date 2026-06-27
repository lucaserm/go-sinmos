package enrollments

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
	enrollmentsRouter := chi.NewRouter()
	enrollmentsRouter.Post("/", h.createEnrollment)
	enrollmentsRouter.Get("/", h.getEnrollments)
	enrollmentsRouter.Get("/{id}", h.getEnrollmentByID)
	enrollmentsRouter.Put("/{id}", h.updateEnrollment)
	enrollmentsRouter.Delete("/{id}", h.deleteEnrollment)

	router.Mount("/enrollments", enrollmentsRouter)
}

func (h *handler) createEnrollment(w http.ResponseWriter, r *http.Request) {
	var payload CreateEnrollmentPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.createEnrollment(r.Context(), payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"enrollment": response,
	})
}

func (h *handler) getEnrollments(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	// Convert offset and limit to int32
	offsetInt := json.ParseInt32(offset, 0)
	limitInt := json.ParseInt32(limit, 10)

	response, err := h.service.getEnrollments(r.Context(), offsetInt, limitInt)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"enrollments": response,
	})
}

func (h *handler) getEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	response, err := h.service.getEnrollmentByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrEnrollmentNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"enrollment": response,
	})
}

func (h *handler) updateEnrollment(w http.ResponseWriter, r *http.Request) {
	var payload UpdateEnrollmentPayload

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
	response, err := h.service.updateEnrollment(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrEnrollmentNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"enrollment": response,
	})
}

func (h *handler) deleteEnrollment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	err := h.service.deleteEnrollment(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrEnrollmentNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "enrollment deleted successfully",
	})
}

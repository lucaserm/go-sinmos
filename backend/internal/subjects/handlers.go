package subjects

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
	subjectsRouter := chi.NewRouter()
	subjectsRouter.Post("/", h.createSubject)
	subjectsRouter.Get("/", h.getSubjects)
	subjectsRouter.Get("/{id}", h.getSubjectByID)
	subjectsRouter.Put("/{id}", h.updateSubject)
	subjectsRouter.Delete("/{id}", h.deleteSubject)

	router.Mount("/subjects", subjectsRouter)
}

func (h *handler) createSubject(w http.ResponseWriter, r *http.Request) {
	var payload CreateSubjectPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.createSubject(r.Context(), payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"subject": response,
	})
}

func (h *handler) getSubjects(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	// Convert offset and limit to int32
	offsetInt := json.ParseInt32(offset, 0)
	limitInt := json.ParseInt32(limit, 10)

	response, err := h.service.getSubjects(r.Context(), offsetInt, limitInt)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"subjects": response,
	})
}

func (h *handler) getSubjectByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	response, err := h.service.getSubjectByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrSubjectNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"subject": response,
	})
}

func (h *handler) updateSubject(w http.ResponseWriter, r *http.Request) {
	var payload UpdateSubjectPayload

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
	response, err := h.service.updateSubject(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrSubjectNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"subject": response,
	})
}

func (h *handler) deleteSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	err := h.service.deleteSubject(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrSubjectNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "subject deleted successfully",
	})
}

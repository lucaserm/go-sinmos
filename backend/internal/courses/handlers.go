package courses

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
	coursesRouter := chi.NewRouter()
	coursesRouter.Post("/", h.createCourse)
	coursesRouter.Get("/", h.getCourses)
	coursesRouter.Get("/{id}", h.getCourseByID)
	coursesRouter.Put("/{id}", h.updateCourse)
	coursesRouter.Delete("/{id}", h.deleteCourse)

	router.Mount("/courses", coursesRouter)
}

func (h *handler) createCourse(w http.ResponseWriter, r *http.Request) {
	var payload CreateCoursePayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.createCourse(r.Context(), payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"course": response,
	})
}

func (h *handler) getCourses(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	// Convert offset and limit to int32
	offsetInt := json.ParseInt32(offset, 0)
	limitInt := json.ParseInt32(limit, 10)

	response, err := h.service.getCourses(r.Context(), offsetInt, limitInt)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"courses": response,
	})
}

func (h *handler) getCourseByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	response, err := h.service.getCourseByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrCourseNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"course": response,
	})
}

func (h *handler) updateCourse(w http.ResponseWriter, r *http.Request) {
	var payload UpdateCoursePayload

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
	response, err := h.service.updateCourse(r.Context(), id, payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrCourseNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"course": response,
	})
}

func (h *handler) deleteCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	err := h.service.deleteCourse(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrCourseNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "course deleted successfully",
	})
}

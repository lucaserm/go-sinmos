package students

import (
	"errors"
	"fmt"
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
	studentRouter := chi.NewRouter()
	studentRouter.Get("/", h.getAllStudents)
	studentRouter.Get("/{id}", h.getStudentByID)
	studentRouter.Post("/", h.createStudent)
	studentRouter.Put("/{id}", h.updateStudent)
	studentRouter.Delete("/{id}", h.deleteStudent)

	router.Mount("/students", studentRouter)
}

func (h *handler) getAllStudents(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	offsetInt := json.ParseInt32(offset, 0)
	limitInt := json.ParseInt32(limit, 10)

	students, err := h.service.getAllStudents(r.Context(), offsetInt, limitInt)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"students": students,
	})
}

func (h *handler) getStudentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	student, err := h.service.getStudentByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrStudentNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, map[string]any{
		"student": student,
	})
}

func (h *handler) createStudent(w http.ResponseWriter, r *http.Request) {
	var student StudentRequest

	if err := json.Read(r, &student); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validation
	if err := json.Validate.Struct(student); err != nil {
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", err))
		return
	}

	response, err := h.service.createStudent(r.Context(), student)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"student": response,
	})
}

func (h *handler) updateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	var student StudentUpdateRequest

	if err := json.Read(r, &student); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validation
	if err := json.Validate.Struct(student); err != nil {
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", err))
		return
	}

	response, err := h.service.updateStudent(r.Context(), id, student)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrStudentNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"student": response,
	})
}

func (h *handler) deleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !json.ValidUUID(w, id) {
		return
	}
	if err := h.service.deleteStudent(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, ErrStudentNotFound)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "student deleted successfully",
	})
}

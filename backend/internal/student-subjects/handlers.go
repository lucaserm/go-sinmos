package studentsubjects

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	studentsubjects := chi.NewRouter()
	studentsubjects.Post("/", h.addStudentSubject)
	studentsubjects.Delete("/", h.removeStudentSubject)

	router.Mount("/students-subjects", studentsubjects)
}

func (h *handler) addStudentSubject(w http.ResponseWriter, r *http.Request) {
	var payload AddStudentSubjectPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.addStudentSubject(r.Context(), payload.StudentID, payload.SubjectID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "student subject added successfully",
	})
}

func (h *handler) removeStudentSubject(w http.ResponseWriter, r *http.Request) {
	var payload RemoveStudentSubjectPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.removeStudentSubject(r.Context(), payload.StudentID, payload.SubjectID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "student subject removed successfully",
	})
}

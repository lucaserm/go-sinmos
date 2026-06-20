package studentguardians

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
	studentguardians := chi.NewRouter()
	studentguardians.Post("/", h.addStudentGuardian)
	studentguardians.Delete("/", h.removeStudentGuardian)

	router.Mount("/students-guardians", studentguardians)
}

func (h *handler) addStudentGuardian(w http.ResponseWriter, r *http.Request) {
	var payload AddStudentGuardianPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.addStudentGuardian(r.Context(), payload.StudentID, payload.GuardianID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "student guardian added successfully",
	})
}

func (h *handler) removeStudentGuardian(w http.ResponseWriter, r *http.Request) {
	var payload RemoveStudentGuardianPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.service.removeStudentGuardian(r.Context(), payload.StudentID, payload.GuardianID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "student guardian removed successfully",
	})
}

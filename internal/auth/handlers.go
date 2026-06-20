package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
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

func (h *handler) RegisterRoutes(router *chi.Mux, repo *repo.Queries) {
	authRouter := chi.NewRouter()
	authRouter.Post("/register", h.register)
	authRouter.Post("/login", h.login)

	authMiddleware := NewMiddleware(repo)
	authRouter.Post("/refresh", h.refreshToken)
	authRouter.Get("/me", authMiddleware.RequiresAuth(h.me))

	router.Mount("/auth", authRouter)
}

func (h *handler) me(w http.ResponseWriter, r *http.Request, user repo.User) {
	json.WriteJSON(w, http.StatusOK, map[string]any{
		"user": UserResponse{
			ID:   user.ID.Bytes,
			Name: user.Name,
			Code: user.Code,
			Role: string(user.Role),
		},
	})
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validation
	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %s", err))
		return
	}

	response, err := h.service.register(r.Context(), payload)

	if err != nil {
		if err == ErrCodeConflict {
			json.WriteError(w, http.StatusConflict, err)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, response)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.login(r.Context(), payload)

	if err != nil {
		if err == ErrInvalidCredentials {
			json.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, response)
}

func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var payload RefreshTokenPayload

	if err := json.Read(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	response, err := h.service.refreshToken(r.Context(), payload)

	if err != nil {
		if err == ErrRefreshTokenExpired {
			json.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		if err == ErrUserNotFound {
			json.WriteError(w, http.StatusNotFound, err)
			return
		}
		if err == ErrRefreshTokenNotFound {
			json.WriteError(w, http.StatusNotFound, err)
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, response)
}

package auth

import (
	"context"

	"github.com/google/uuid"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
)

type Role string

const (
	ADMIN     Role = "ADMIN"
	SUPPORT   Role = "SUPPORT"
	RECEPTION Role = "RECEPTION"
)

func (r Role) SQLC() repo.Role {
	return repo.Role(r)
}

type (
	RegisterPayload struct {
		Name     string `json:"name" validate:"required,min=3"`
		Code     string `json:"code" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
		Role     Role   `json:"role" validate:"required,oneof=ADMIN SUPPORT RECEPTION"`
	}

	LoginPayload struct {
		Code     string `json:"code" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	RefreshTokenPayload struct {
		RefreshToken string `json:"refreshToken" validate:"required,uuid"`
	}
)

type (
	UserResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
		Code string    `json:"code"`
		Role string    `json:"role"`
	}

	AuthResponse struct {
		User         UserResponse `json:"user"`
		AccessToken  string       `json:"accessToken"`
		RefreshToken string       `json:"refreshToken"`
	}

	RefreshTokenResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)

type Service interface {
	register(ctx context.Context, payload RegisterPayload) (AuthResponse, error)
	login(ctx context.Context, payload LoginPayload) (AuthResponse, error)
	refreshToken(ctx context.Context, payload RefreshTokenPayload) (RefreshTokenResponse, error)
}

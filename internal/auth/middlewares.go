package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/go-sinmos/internal/json"
	"github.com/lucaserm/go-sinmos/internal/jwt"
)

type authedHandler func(http.ResponseWriter, *http.Request, repo.User)

type middleware struct {
	repo *repo.Queries
}

func NewMiddleware(repo *repo.Queries) *middleware {
	return &middleware{
		repo: repo,
	}
}

func (m *middleware) RequiresAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetJWTToken(r.Header)
		if err != nil {
			json.WriteError(w, http.StatusUnauthorized, fmt.Errorf("auth error: %v", err))
			return
		}

		userIdStr, err := jwt.ValidateToken(token, true)
		if err != nil {
			json.WriteError(w, http.StatusUnauthorized, fmt.Errorf("auth error: %v", err))
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			json.WriteError(w, http.StatusUnauthorized, fmt.Errorf("auth error: %v", err))
			return
		}

		user, err := m.repo.GetUserById(r.Context(), pgtype.UUID{Bytes: userId, Valid: true})
		if err != nil {
			if err == pgx.ErrNoRows {
				json.WriteError(w, http.StatusUnauthorized, errors.New("user not found"))
				return
			}
			json.WriteError(w, http.StatusInternalServerError, errors.New("internal error"))
			return
		}

		handler(w, r, user)
	}
}

func GetJWTToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authorization header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of authorization header")
	}

	return vals[1], nil
}

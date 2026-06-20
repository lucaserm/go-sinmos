package auth

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/go-sinmos/internal/jwt"
	"github.com/lucaserm/go-sinmos/internal/utils"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (svc *svc) register(ctx context.Context, payload RegisterPayload) (AuthResponse, error) {
	user, err := svc.repo.GetUserByCode(ctx, payload.Code)
	if err != nil && err != pgx.ErrNoRows {
		return AuthResponse{}, err
	}

	if err == nil {
		return AuthResponse{}, ErrCodeConflict
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return AuthResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err = svc.repo.CreateUser(ctx, repo.CreateUserParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		Name:           payload.Name,
		Code:           payload.Code,
		HashedPassword: string(hashedPassword),
		Role:           payload.Role.SQLC(),
	})

	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, err := svc.repo.CreateRefreshToken(
		ctx,
		repo.CreateRefreshTokenParams{
			Token: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			UserID: user.ID,
			ExpiresAt: pgtype.Timestamptz{
				Time:  time.Now().Add(7 * 24 * time.Hour),
				Valid: true,
			},
		},
	)
	if err != nil {
		return AuthResponse{}, err
	}

	accessToken, err := jwt.CreateToken(user.ID.String())
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		User: UserResponse{
			ID:   user.ID.Bytes,
			Name: user.Name,
			Code: user.Code,
			Role: string(user.Role),
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token.String(),
	}, nil
}

func (svc *svc) login(ctx context.Context, payload LoginPayload) (AuthResponse, error) {
	user, err := svc.repo.GetUserByCode(ctx, payload.Code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	isValid := utils.CheckPassword(user.HashedPassword, payload.Password)
	if !isValid {
		return AuthResponse{}, ErrInvalidCredentials
	}

	refreshToken, err := svc.repo.GetRefreshTokenByUserId(ctx, repo.GetRefreshTokenByUserIdParams{
		UserID: user.ID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		if err != pgx.ErrNoRows {
			return AuthResponse{}, err
		}
		refreshToken, err = svc.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
			Token:     pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:    user.ID,
			ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
		})
		if err != nil {
			return AuthResponse{}, err
		}
	}

	accessToken, err := jwt.CreateToken(user.ID.String())
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		User: UserResponse{
			ID:   user.ID.Bytes,
			Name: user.Name,
			Code: user.Code,
			Role: string(user.Role),
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token.String(),
	}, nil
}

func (svc *svc) refreshToken(ctx context.Context, payload RefreshTokenPayload) (RefreshTokenResponse, error) {
	refreshToken, err := svc.repo.GetRefreshTokenByToken(ctx, repo.GetRefreshTokenByTokenParams{
		Token: pgtype.UUID{
			Bytes: uuid.MustParse(payload.RefreshToken),
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return RefreshTokenResponse{}, ErrRefreshTokenExpired
		}
		return RefreshTokenResponse{}, err
	}

	if refreshToken.Token.String() != payload.RefreshToken {
		return RefreshTokenResponse{}, ErrRefreshTokenNotFound
	}

	token, err := jwt.CreateToken(refreshToken.UserID.String())
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	err = svc.repo.DeleteRefreshTokenByUserId(ctx, refreshToken.UserID)

	if err != nil {
		return RefreshTokenResponse{}, err
	}

	refreshToken, err = svc.repo.CreateRefreshToken(ctx, repo.CreateRefreshTokenParams{
		Token: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		UserID: refreshToken.UserID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(7 * 24 * time.Hour),
			Valid: true,
		},
	})

	if err != nil {
		return RefreshTokenResponse{}, err
	}

	return RefreshTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken.Token.String(),
	}, nil
}

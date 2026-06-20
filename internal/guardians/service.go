package guardians

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) createGuardian(ctx context.Context, payload CreateGuardianPayload) (GuardianResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return GuardianResponse{}, err
	}

	guardian, err := s.repo.CreateGuardian(ctx, repo.CreateGuardianParams{
		ID:    pgtype.UUID{Bytes: uid, Valid: true},
		Name:  payload.Name,
		Email: payload.Email,
		Phone: payload.Phone,
	})
	if err != nil {
		return GuardianResponse{}, err
	}

	return guardianToResponse(guardian), nil
}

func (s *svc) getGuardians(ctx context.Context, offset, limit int32) ([]GuardianResponse, error) {
	guardians, err := s.repo.GetGuardians(ctx, repo.GetGuardiansParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return guardiansToResponse(guardians), nil
}

func (s *svc) getGuardianByID(ctx context.Context, id string) (GuardianResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return GuardianResponse{}, err
	}
	guardian, err := s.repo.GetGuardianByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return GuardianResponse{}, err
	}

	return guardianToResponse(guardian), nil
}

func (s *svc) updateGuardian(ctx context.Context, id string, payload UpdateGuardianPayload) (GuardianResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return GuardianResponse{}, err
	}
	guardian, err := s.repo.UpdateGuardian(ctx, repo.UpdateGuardianParams{
		ID:    pgtype.UUID{Bytes: uidParsed, Valid: true},
		Name:  payload.Name,
		Email: payload.Email,
		Phone: payload.Phone,
	})
	if err != nil {
		return GuardianResponse{}, err
	}

	return guardianToResponse(guardian), nil
}

func (s *svc) deleteGuardian(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteGuardian(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func guardianToResponse(guardian repo.Guardian) GuardianResponse {
	return GuardianResponse{
		ID:    guardian.ID.String(),
		Name:  guardian.Name,
		Email: guardian.Email,
		Phone: guardian.Phone,
	}
}

func guardiansToResponse(guardians []repo.Guardian) []GuardianResponse {
	var response []GuardianResponse = make([]GuardianResponse, 0, len(guardians))
	for _, guardian := range guardians {
		response = append(response, guardianToResponse(guardian))
	}
	return response
}

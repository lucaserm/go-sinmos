package warnings

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

func (s *svc) createWarning(ctx context.Context, payload CreateWarningPayload) (WarningResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return WarningResponse{}, err
	}

	warning, err := s.repo.CreateWarning(ctx, repo.CreateWarningParams{
		ID:           pgtype.UUID{Bytes: uid, Valid: true},
		OccurrenceID: pgtype.UUID{Bytes: uuid.MustParse(payload.OccurrenceID), Valid: true},
		Report:       payload.Report,
	})
	if err != nil {
		return WarningResponse{}, err
	}

	return warningToResponse(warning), nil
}

func (s *svc) getWarnings(ctx context.Context, offset, limit int32) ([]WarningResponse, error) {
	warnings, err := s.repo.GetAllWarnings(ctx, repo.GetAllWarningsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return warningsToResponse(warnings), nil
}

func (s *svc) getWarningByID(ctx context.Context, id string) (WarningResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return WarningResponse{}, err
	}
	warning, err := s.repo.GetWarningByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return WarningResponse{}, err
	}

	return warningToResponse(warning), nil
}

func (s *svc) updateWarning(ctx context.Context, id string, payload UpdateWarningPayload) (WarningResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return WarningResponse{}, err
	}

	// Read-modify-write: both fields are optional. Parse occurrenceId only when
	// provided (uuid.MustParse would panic on an empty value).
	existing, err := s.repo.GetWarningByID(ctx, pgtype.UUID{Bytes: uidParsed, Valid: true})
	if err != nil {
		return WarningResponse{}, err
	}

	params := repo.UpdateWarningParams{
		ID:           existing.ID,
		OccurrenceID: existing.OccurrenceID,
		Report:       existing.Report,
	}
	if payload.OccurrenceID != "" {
		occurrenceID, err := uuid.Parse(payload.OccurrenceID)
		if err != nil {
			return WarningResponse{}, err
		}
		params.OccurrenceID = pgtype.UUID{Bytes: occurrenceID, Valid: true}
	}
	if payload.Report != "" {
		params.Report = payload.Report
	}

	warning, err := s.repo.UpdateWarning(ctx, params)
	if err != nil {
		return WarningResponse{}, err
	}

	return warningToResponse(warning), nil
}

func (s *svc) deleteWarning(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteWarning(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func warningToResponse(warning repo.Warning) WarningResponse {
	return WarningResponse{
		ID:           warning.ID.String(),
		OccurrenceID: warning.OccurrenceID.String(),
		Report:       warning.Report,
	}
}

func warningsToResponse(warnings []repo.Warning) []WarningResponse {
	var response []WarningResponse = make([]WarningResponse, 0, len(warnings))
	for _, warning := range warnings {
		response = append(response, warningToResponse(warning))
	}
	return response
}

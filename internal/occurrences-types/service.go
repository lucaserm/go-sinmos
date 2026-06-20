package occurrencestypes

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

func (s *svc) createOccurrenceType(ctx context.Context, payload CreateOccurrenceTypePayload) (OccurrenceTypeResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return OccurrenceTypeResponse{}, err
	}

	occurrenceType, err := s.repo.CreateOccurrenceType(ctx, repo.CreateOccurrenceTypeParams{
		ID:          pgtype.UUID{Bytes: uid, Valid: true},
		Code:        payload.Code,
		Description: payload.Description,
		Severity:    payload.Severity,
	})
	if err != nil {
		return OccurrenceTypeResponse{}, err
	}

	return occurrenceTypeToResponse(occurrenceType), nil
}

func (s *svc) getOccurrenceTypes(ctx context.Context, offset, limit int32) ([]OccurrenceTypeResponse, error) {
	occurrenceTypes, err := s.repo.GetAllOccurrenceTypes(ctx, repo.GetAllOccurrenceTypesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return occurrenceTypesToResponse(occurrenceTypes), nil
}

func (s *svc) getOccurrenceTypeByID(ctx context.Context, id string) (OccurrenceTypeResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return OccurrenceTypeResponse{}, err
	}
	occurrenceType, err := s.repo.GetOccurrenceTypeByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return OccurrenceTypeResponse{}, err
	}

	return occurrenceTypeToResponse(occurrenceType), nil
}

func (s *svc) updateOccurrenceType(ctx context.Context, id string, payload UpdateOccurrenceTypePayload) (OccurrenceTypeResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return OccurrenceTypeResponse{}, err
	}

	occurrenceType, err := s.repo.UpdateOccurrenceType(ctx, repo.UpdateOccurrenceTypeParams{
		ID:          pgtype.UUID{Bytes: uidParsed, Valid: true},
		Code:        payload.Code,
		Description: payload.Description,
		Severity:    payload.Severity,
	})
	if err != nil {
		return OccurrenceTypeResponse{}, err
	}

	return occurrenceTypeToResponse(occurrenceType), nil
}

func (s *svc) deleteOccurrenceType(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteOccurrenceType(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func occurrenceTypeToResponse(occurrenceType repo.OccurrenceType) OccurrenceTypeResponse {
	return OccurrenceTypeResponse{
		ID:          occurrenceType.ID.String(),
		Code:        occurrenceType.Code,
		Description: occurrenceType.Description,
		Severity:    occurrenceType.Severity,
	}
}

func occurrenceTypesToResponse(occurrenceTypes []repo.OccurrenceType) []OccurrenceTypeResponse {
	var response []OccurrenceTypeResponse = make([]OccurrenceTypeResponse, 0, len(occurrenceTypes))
	for _, occurrenceType := range occurrenceTypes {
		response = append(response, occurrenceTypeToResponse(occurrenceType))
	}
	return response
}

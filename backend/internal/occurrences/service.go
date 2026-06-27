package occurrences

import (
	"context"
	"time"

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

func (s *svc) createOccurrence(ctx context.Context, payload CreateOccurrencePayload, userId uuid.UUID) (OccurrenceResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return OccurrenceResponse{}, err
	}

	occurrenceTypeID, err := uuid.Parse(payload.OccurrenceTypeID)
	if err != nil {
		return OccurrenceResponse{}, err
	}

	studentID, err := uuid.Parse(payload.StudentID)
	if err != nil {
		return OccurrenceResponse{}, err
	}

	occurredAt, err := time.Parse(time.RFC3339, payload.OccurredAt)
	if err != nil {
		return OccurrenceResponse{}, err
	}

	// userRelatedId is optional — leave the column NULL when omitted.
	userRelatedID := pgtype.UUID{}
	if payload.UserRelatedID != "" {
		parsed, err := uuid.Parse(payload.UserRelatedID)
		if err != nil {
			return OccurrenceResponse{}, err
		}
		userRelatedID = pgtype.UUID{Bytes: parsed, Valid: true}
	}

	// status is optional — default to PENDING (matches the column default).
	status := payload.Status
	if status == "" {
		status = string(repo.OccurrenceStatusPENDING)
	}

	occurrence, err := s.repo.CreateOccurrence(ctx, repo.CreateOccurrenceParams{
		ID:               pgtype.UUID{Bytes: uid, Valid: true},
		UserID:           pgtype.UUID{Bytes: userId, Valid: true},
		OccurrenceTypeID: pgtype.UUID{Bytes: occurrenceTypeID, Valid: true},
		StudentID:        pgtype.UUID{Bytes: studentID, Valid: true},
		OccurredAt:       pgtype.Timestamptz{Time: occurredAt, Valid: true},
		UserRelatedID:    userRelatedID,
		Status:           repo.OccurrenceStatus(status),
	})
	if err != nil {
		return OccurrenceResponse{}, err
	}

	return occurrenceToResponse(occurrence), nil
}

func (s *svc) getOccurrences(ctx context.Context, offset, limit int32) ([]OccurrenceResponse, error) {
	occurrences, err := s.repo.GetAllOccurrences(ctx, repo.GetAllOccurrencesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return occurrencesToResponse(occurrences), nil
}

func (s *svc) getOccurrenceByID(ctx context.Context, id string) (OccurrenceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return OccurrenceResponse{}, err
	}
	occurrence, err := s.repo.GetOccurrenceByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return OccurrenceResponse{}, err
	}

	return occurrenceToResponse(occurrence), nil
}

func (s *svc) updateOccurrence(ctx context.Context, id string, payload UpdateOccurrencePayload) (OccurrenceResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return OccurrenceResponse{}, err
	}

	// Read-modify-write: every field on UpdateOccurrencePayload is optional, so
	// start from the existing row and only override what the caller provided.
	// This keeps a partial update (e.g. status only) from blanking out the rest.
	existing, err := s.repo.GetOccurrenceByID(ctx, pgtype.UUID{Bytes: uidParsed, Valid: true})
	if err != nil {
		return OccurrenceResponse{}, err
	}

	params := repo.UpdateOccurrenceParams{
		ID:               existing.ID,
		UserID:           existing.UserID,
		OccurrenceTypeID: existing.OccurrenceTypeID,
		StudentID:        existing.StudentID,
		OccurredAt:       existing.OccurredAt,
		UserRelatedID:    existing.UserRelatedID,
		Status:           existing.Status,
	}

	if payload.UserID != "" {
		params.UserID = pgtype.UUID{Bytes: uuid.MustParse(payload.UserID), Valid: true}
	}
	if payload.OccurrenceTypeID != "" {
		params.OccurrenceTypeID = pgtype.UUID{Bytes: uuid.MustParse(payload.OccurrenceTypeID), Valid: true}
	}
	if payload.StudentID != "" {
		params.StudentID = pgtype.UUID{Bytes: uuid.MustParse(payload.StudentID), Valid: true}
	}
	if payload.OccurredAt != "" {
		occurredAt, err := time.Parse(time.RFC3339, payload.OccurredAt)
		if err != nil {
			return OccurrenceResponse{}, err
		}
		params.OccurredAt = pgtype.Timestamptz{Time: occurredAt, Valid: true}
	}
	if payload.UserRelatedID != "" {
		params.UserRelatedID = pgtype.UUID{Bytes: uuid.MustParse(payload.UserRelatedID), Valid: true}
	}
	if payload.Status != "" {
		params.Status = repo.OccurrenceStatus(payload.Status)
	}

	occurrence, err := s.repo.UpdateOccurrence(ctx, params)
	if err != nil {
		return OccurrenceResponse{}, err
	}

	return occurrenceToResponse(occurrence), nil
}

func (s *svc) deleteOccurrence(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteOccurrence(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func occurrenceToResponse(occurrence repo.Occurrence) OccurrenceResponse {
	return OccurrenceResponse{
		ID:               occurrence.ID.String(),
		UserID:           occurrence.UserID.String(),
		OccurrenceTypeID: occurrence.OccurrenceTypeID.String(),
		StudentID:        occurrence.StudentID.String(),
		UserRelatedID:    occurrence.UserRelatedID.String(),
		Status:           string(occurrence.Status),
		OccurredAt:       occurrence.OccurredAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func occurrencesToResponse(occurrences []repo.Occurrence) []OccurrenceResponse {
	var response []OccurrenceResponse = make([]OccurrenceResponse, 0, len(occurrences))
	for _, occurrence := range occurrences {
		response = append(response, occurrenceToResponse(occurrence))
	}
	return response
}

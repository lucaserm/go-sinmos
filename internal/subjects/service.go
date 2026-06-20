package subjects

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

func (s *svc) createSubject(ctx context.Context, payload CreateSubjectPayload) (SubjectResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return SubjectResponse{}, err
	}

	subject, err := s.repo.CreateSubject(ctx, repo.CreateSubjectParams{
		ID:       pgtype.UUID{Bytes: uid, Valid: true},
		Name:     payload.Name,
		CourseID: pgtype.UUID{Bytes: uuid.MustParse(payload.CourseID), Valid: true},
		Semester: payload.Semester,
		Section:  payload.Section,
	})
	if err != nil {
		return SubjectResponse{}, err
	}

	return subjectToResponse(subject), nil
}

func (s *svc) getSubjects(ctx context.Context, offset, limit int32) ([]SubjectResponse, error) {
	subjects, err := s.repo.GetAllSubjects(ctx, repo.GetAllSubjectsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return subjectsToResponse(subjects), nil
}

func (s *svc) getSubjectByID(ctx context.Context, id string) (SubjectResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return SubjectResponse{}, err
	}
	subject, err := s.repo.GetSubjectByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return SubjectResponse{}, err
	}

	return subjectToResponse(subject), nil
}

func (s *svc) updateSubject(ctx context.Context, id string, payload UpdateSubjectPayload) (SubjectResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return SubjectResponse{}, err
	}

	subject, err := s.repo.UpdateSubject(ctx, repo.UpdateSubjectParams{
		ID:       pgtype.UUID{Bytes: uidParsed, Valid: true},
		Name:     payload.Name,
		Semester: payload.Semester,
		Section:  payload.Section,
	})
	if err != nil {
		return SubjectResponse{}, err
	}

	return subjectToResponse(subject), nil
}

func (s *svc) deleteSubject(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteSubject(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func subjectToResponse(subject repo.Subject) SubjectResponse {
	return SubjectResponse{
		ID:       subject.ID.String(),
		Name:     subject.Name,
		Semester: subject.Semester,
		Section:  subject.Section,
		CourseID: subject.CourseID.String(),
	}
}

func subjectsToResponse(subjects []repo.Subject) []SubjectResponse {
	var response []SubjectResponse = make([]SubjectResponse, 0, len(subjects))
	for _, subject := range subjects {
		response = append(response, subjectToResponse(subject))
	}
	return response
}

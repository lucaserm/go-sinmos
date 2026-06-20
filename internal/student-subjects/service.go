package studentsubjects

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

func (s *svc) addStudentSubject(ctx context.Context, studentID, subjectID string) error {
	studentUid, err := uuid.Parse(studentID)
	if err != nil {
		return err
	}

	subjectUid, err := uuid.Parse(subjectID)
	if err != nil {
		return err
	}
	return s.repo.AddStudentSubject(ctx, repo.AddStudentSubjectParams{
		StudentID: pgtype.UUID{Bytes: studentUid, Valid: true},
		SubjectID: pgtype.UUID{Bytes: subjectUid, Valid: true},
	})
}

func (s *svc) removeStudentSubject(ctx context.Context, studentID, subjectID string) error {
	studentUid, err := uuid.Parse(studentID)
	if err != nil {
		return err
	}

	subjectUid, err := uuid.Parse(subjectID)
	if err != nil {
		return err
	}

	return s.repo.RemoveStudentSubject(ctx, repo.RemoveStudentSubjectParams{
		StudentID: pgtype.UUID{Bytes: studentUid, Valid: true},
		SubjectID: pgtype.UUID{Bytes: subjectUid, Valid: true},
	})
}

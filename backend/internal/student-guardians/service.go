package studentguardians

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

func (s *svc) addStudentGuardian(ctx context.Context, studentID, guardianID string) error {
	studentUid, err := uuid.Parse(studentID)
	if err != nil {
		return err
	}

	guardianUid, err := uuid.Parse(guardianID)
	if err != nil {
		return err
	}
	return s.repo.AddStudentGuardian(ctx, repo.AddStudentGuardianParams{
		StudentID:  pgtype.UUID{Bytes: studentUid, Valid: true},
		GuardianID: pgtype.UUID{Bytes: guardianUid, Valid: true},
	})
}

func (s *svc) removeStudentGuardian(ctx context.Context, studentID, guardianID string) error {
	studentUid, err := uuid.Parse(studentID)
	if err != nil {
		return err
	}

	guardianUid, err := uuid.Parse(guardianID)
	if err != nil {
		return err
	}

	return s.repo.RemoveStudentGuardian(ctx, repo.RemoveStudentGuardianParams{
		StudentID:  pgtype.UUID{Bytes: studentUid, Valid: true},
		GuardianID: pgtype.UUID{Bytes: guardianUid, Valid: true},
	})
}

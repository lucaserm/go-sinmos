package enrollments

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

func (s *svc) createEnrollment(ctx context.Context, payload CreateEnrollmentPayload) (EnrollmentResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return EnrollmentResponse{}, err
	}

	enrollment, err := s.repo.CreateEnrollment(ctx, repo.CreateEnrollmentParams{
		ID:        pgtype.UUID{Bytes: uid, Valid: true},
		StudentID: pgtype.UUID{Bytes: uuid.MustParse(payload.StudentID), Valid: true},
		CourseID:  pgtype.UUID{Bytes: uuid.MustParse(payload.CourseID), Valid: true},
		Year:      payload.Year,
	})
	if err != nil {
		return EnrollmentResponse{}, err
	}

	return enrollmentToResponse(enrollment), nil
}

func (s *svc) getEnrollments(ctx context.Context, offset, limit int32) ([]EnrollmentResponse, error) {
	enrollments, err := s.repo.GetAllEnrollments(ctx, repo.GetAllEnrollmentsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return enrollmentsToResponse(enrollments), nil
}

func (s *svc) getEnrollmentByID(ctx context.Context, id string) (EnrollmentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return EnrollmentResponse{}, err
	}
	enrollment, err := s.repo.GetEnrollmentByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return EnrollmentResponse{}, err
	}

	return enrollmentToResponse(enrollment), nil
}

func (s *svc) updateEnrollment(ctx context.Context, id string, payload UpdateEnrollmentPayload) (EnrollmentResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return EnrollmentResponse{}, err
	}

	// Read-modify-write: all fields are optional. Parse the UUIDs only when
	// provided (uuid.MustParse would panic on empty values).
	existing, err := s.repo.GetEnrollmentByID(ctx, pgtype.UUID{Bytes: uidParsed, Valid: true})
	if err != nil {
		return EnrollmentResponse{}, err
	}

	params := repo.UpdateEnrollmentParams{
		ID:        existing.ID,
		StudentID: existing.StudentID,
		CourseID:  existing.CourseID,
		Year:      existing.Year,
	}
	if payload.StudentID != "" {
		studentID, err := uuid.Parse(payload.StudentID)
		if err != nil {
			return EnrollmentResponse{}, err
		}
		params.StudentID = pgtype.UUID{Bytes: studentID, Valid: true}
	}
	if payload.CourseID != "" {
		courseID, err := uuid.Parse(payload.CourseID)
		if err != nil {
			return EnrollmentResponse{}, err
		}
		params.CourseID = pgtype.UUID{Bytes: courseID, Valid: true}
	}
	if payload.Year != 0 {
		params.Year = payload.Year
	}

	enrollment, err := s.repo.UpdateEnrollment(ctx, params)
	if err != nil {
		return EnrollmentResponse{}, err
	}

	return enrollmentToResponse(enrollment), nil
}

func (s *svc) deleteEnrollment(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteEnrollment(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func enrollmentToResponse(enrollment repo.Enrollment) EnrollmentResponse {
	return EnrollmentResponse{
		ID:        enrollment.ID.String(),
		StudentID: enrollment.StudentID.String(),
		CourseID:  enrollment.CourseID.String(),
		Year:      enrollment.Year,
	}
}

func enrollmentsToResponse(enrollments []repo.Enrollment) []EnrollmentResponse {
	var response []EnrollmentResponse = make([]EnrollmentResponse, 0, len(enrollments))
	for _, enrollment := range enrollments {
		response = append(response, enrollmentToResponse(enrollment))
	}
	return response
}

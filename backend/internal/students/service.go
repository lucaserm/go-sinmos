package students

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

func (s *svc) getAllStudents(ctx context.Context, offset, limit int32) ([]StudentResponse, error) {
	students, err := s.repo.GetAllStudents(ctx, repo.GetAllStudentsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []StudentResponse{}, err
	}
	return studentsToResponse(students), nil
}

func (s *svc) getStudentByID(ctx context.Context, id string) (StudentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return StudentResponse{}, err
	}
	student, err := s.repo.GetStudentByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return StudentResponse{}, err
	}
	return studentToResponse(student), nil
}

func (s *svc) createStudent(ctx context.Context, student StudentRequest) (StudentResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return StudentResponse{}, err
	}

	studentDB, err := s.repo.CreateStudent(ctx, repo.CreateStudentParams{
		ID:       pgtype.UUID{Bytes: id, Valid: true},
		Cpf:      student.CPF,
		Ra:       student.RA,
		PhotoUrl: pgtype.Text{String: student.PhotoURL, Valid: true},
		Name:     student.Name,
		Email:    student.Email,
	})

	if err != nil {
		return StudentResponse{}, err
	}
	return studentToResponse(studentDB), nil
}

func (s *svc) updateStudent(ctx context.Context, id string, student StudentUpdateRequest) (StudentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return StudentResponse{}, err
	}

	// Read-modify-write: every payload field is optional, so override only the
	// fields that were provided rather than blanking the rest.
	existing, err := s.repo.GetStudentByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return StudentResponse{}, err
	}

	params := repo.UpdateStudentParams{
		ID:       existing.ID,
		Cpf:      existing.Cpf,
		Ra:       existing.Ra,
		PhotoUrl: existing.PhotoUrl,
		Name:     existing.Name,
		Email:    existing.Email,
	}
	if student.CPF != "" {
		params.Cpf = student.CPF
	}
	if student.RA != "" {
		params.Ra = student.RA
	}
	if student.PhotoURL != "" {
		params.PhotoUrl = pgtype.Text{String: student.PhotoURL, Valid: true}
	}
	if student.Name != "" {
		params.Name = student.Name
	}
	if student.Email != "" {
		params.Email = student.Email
	}

	studentDB, err := s.repo.UpdateStudent(ctx, params)
	if err != nil {
		return StudentResponse{}, err
	}
	return studentToResponse(studentDB), nil
}

func (s *svc) deleteStudent(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteStudent(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func studentToResponse(student repo.Student) StudentResponse {
	return StudentResponse{
		ID:       student.ID.String(),
		CPF:      student.Cpf,
		RA:       student.Ra,
		PhotoURL: student.PhotoUrl.String,
		Name:     student.Name,
		Email:    student.Email,
	}
}

func studentsToResponse(students []repo.Student) []StudentResponse {
	responses := make([]StudentResponse, 0, len(students))
	for _, student := range students {
		responses = append(responses, studentToResponse(student))
	}
	return responses
}

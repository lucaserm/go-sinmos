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
	student, err := s.repo.GetStudentByID(ctx, pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true})
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
	studentDB, err := s.repo.UpdateStudent(ctx, repo.UpdateStudentParams{
		ID:       pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true},
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

func (s *svc) deleteStudent(ctx context.Context, id string) error {
	err := s.repo.DeleteStudent(ctx, pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true})
	return err
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

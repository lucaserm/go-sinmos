package students

import "context"

type (
	StudentRequest struct {
		CPF      string `json:"cpf" validate:"required,max=30"`
		RA       string `json:"ra" validate:"required,max=20"`
		PhotoURL string `json:"photo_url"`
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
	}
	StudentUpdateRequest struct {
		CPF      string `json:"cpf,omitempty" validate:"omitempty,max=30"`
		RA       string `json:"ra,omitempty" validate:"omitempty,max=20"`
		PhotoURL string `json:"photo_url,omitempty"`
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
	}
)

type (
	StudentResponse struct {
		ID       string `json:"id"`
		CPF      string `json:"cpf"`
		RA       string `json:"ra"`
		PhotoURL string `json:"photo_url"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
)

type Service interface {
	getAllStudents(ctx context.Context, offset, limit int32) ([]StudentResponse, error)
	getStudentByID(ctx context.Context, id string) (StudentResponse, error)
	createStudent(ctx context.Context, student StudentRequest) (StudentResponse, error)
	updateStudent(ctx context.Context, id string, student StudentUpdateRequest) (StudentResponse, error)
	deleteStudent(ctx context.Context, id string) error
}

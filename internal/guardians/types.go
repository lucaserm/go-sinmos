package guardians

import "context"

type (
	CreateGuardianPayload struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Phone string `json:"phone" validate:"required"`
	}

	UpdateGuardianPayload struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty" validate:"omitempty,email"`
		Phone string `json:"phone,omitempty"`
	}

	GuardianResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
)

type Service interface {
	createGuardian(ctx context.Context, payload CreateGuardianPayload) (GuardianResponse, error)
	getGuardians(ctx context.Context, offset, limit int32) ([]GuardianResponse, error)
	getGuardianByID(ctx context.Context, id string) (GuardianResponse, error)
	updateGuardian(ctx context.Context, id string, payload UpdateGuardianPayload) (GuardianResponse, error)
	deleteGuardian(ctx context.Context, id string) error
}

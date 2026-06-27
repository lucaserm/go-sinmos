package warnings

import "context"

type (
	CreateWarningPayload struct {
		OccurrenceID string `json:"occurrenceId" validate:"required,uuid"`
		Report       string `json:"report" validate:"required"`
	}

	UpdateWarningPayload struct {
		OccurrenceID string `json:"occurrenceId,omitempty" validate:"omitempty,uuid"`
		Report       string `json:"report,omitempty"`
	}

	WarningResponse struct {
		ID           string `json:"id"`
		OccurrenceID string `json:"occurrenceId"`
		Report       string `json:"report"`
	}
)

type Service interface {
	createWarning(ctx context.Context, payload CreateWarningPayload) (WarningResponse, error)
	getWarnings(ctx context.Context, offset, limit int32) ([]WarningResponse, error)
	getWarningByID(ctx context.Context, id string) (WarningResponse, error)
	updateWarning(ctx context.Context, id string, payload UpdateWarningPayload) (WarningResponse, error)
	deleteWarning(ctx context.Context, id string) error
}

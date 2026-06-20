package occurrencestypes

import "context"

type (
	CreateOccurrenceTypePayload struct {
		Code        string `json:"code" validate:"required"`
		Description string `json:"description" validate:"required"`
		Severity    int32  `json:"severity" validate:"required"`
	}

	UpdateOccurrenceTypePayload struct {
		Code        string `json:"code,omitempty"`
		Description string `json:"description,omitempty"`
		Severity    int32  `json:"severity,omitempty"`
	}

	OccurrenceTypeResponse struct {
		ID          string `json:"id"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Severity    int32  `json:"severity"`
	}
)

type Service interface {
	createOccurrenceType(ctx context.Context, payload CreateOccurrenceTypePayload) (OccurrenceTypeResponse, error)
	getOccurrenceTypes(ctx context.Context, offset, limit int32) ([]OccurrenceTypeResponse, error)
	getOccurrenceTypeByID(ctx context.Context, id string) (OccurrenceTypeResponse, error)
	updateOccurrenceType(ctx context.Context, id string, payload UpdateOccurrenceTypePayload) (OccurrenceTypeResponse, error)
	deleteOccurrenceType(ctx context.Context, id string) error
}

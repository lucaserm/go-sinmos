package occurrences

import (
	"context"

	"github.com/google/uuid"
)

type (
	CreateOccurrencePayload struct {
		OccurrenceTypeID string `json:"occurrenceTypeId" validate:"required,uuid"`
		StudentID        string `json:"studentId" validate:"required,uuid"`
		OccurredAt       string `json:"occurredAt" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
		UserRelatedID    string `json:"userRelatedId,omitempty" validate:"omitempty,uuid"`
		Status           string `json:"status,omitempty" validate:"omitempty,oneof=PENDING APPROVED REPROVED"`
	}

	UpdateOccurrencePayload struct {
		UserID           string `json:"userId,omitempty" validate:"omitempty,uuid"`
		OccurrenceTypeID string `json:"occurrenceTypeId,omitempty" validate:"omitempty,uuid"`
		StudentID        string `json:"studentId,omitempty" validate:"omitempty,uuid"`
		OccurredAt       string `json:"occurredAt,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
		UserRelatedID    string `json:"userRelatedId,omitempty" validate:"omitempty,uuid"`
		Status           string `json:"status,omitempty" validate:"omitempty,oneof=PENDING APPROVED REPROVED"`
	}

	OccurrenceResponse struct {
		ID               string `json:"id"`
		UserID           string `json:"userId"`
		OccurrenceTypeID string `json:"occurrenceTypeId"`
		StudentID        string `json:"studentId"`
		OccurredAt       string `json:"occurredAt"`
		UserRelatedID    string `json:"userRelatedId,omitempty"`
		Status           string `json:"status"`
	}
)

type Service interface {
	createOccurrence(ctx context.Context, payload CreateOccurrencePayload, userId uuid.UUID) (OccurrenceResponse, error)
	getOccurrences(ctx context.Context, offset, limit int32) ([]OccurrenceResponse, error)
	getOccurrenceByID(ctx context.Context, id string) (OccurrenceResponse, error)
	updateOccurrence(ctx context.Context, id string, payload UpdateOccurrencePayload) (OccurrenceResponse, error)
	deleteOccurrence(ctx context.Context, id string) error
}

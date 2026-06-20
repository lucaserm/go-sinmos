package permissions

import (
	"context"
)

type (
	CreatePermissionPayload struct {
		StudentID    string `json:"studentId" validate:"required"`
		Type         string `json:"type" validate:"required,oneof=LEAVE STAY"`
		Description  string `json:"description" validate:"required"`
		RequestedAt  int64  `json:"requestedAt" validate:"required"`
		ScheduledFor int64  `json:"scheduledFor" validate:"required"`
	}

	UpdatePermissionPayload struct {
		StudentID    string `json:"studentId,omitempty"`
		Type         string `json:"type,omitempty"`
		Description  string `json:"description,omitempty"`
		RequestedAt  int64  `json:"requestedAt,omitempty"`
		ScheduledFor int64  `json:"scheduledFor,omitempty"`
	}

	PermissionResponse struct {
		ID           string `json:"id"`
		StudentID    string `json:"studentId"`
		Type         string `json:"type"`
		Description  string `json:"description"`
		RequestedAt  string `json:"requestedAt"`
		ScheduledFor string `json:"scheduledFor"`
	}
)

type Service interface {
	createPermission(ctx context.Context, payload CreatePermissionPayload) (PermissionResponse, error)
	getPermissions(ctx context.Context, offset, limit int32) ([]PermissionResponse, error)
	getPermissionByID(ctx context.Context, id string) (PermissionResponse, error)
	updatePermission(ctx context.Context, id string, payload UpdatePermissionPayload) (PermissionResponse, error)
	deletePermission(ctx context.Context, id string) error
}

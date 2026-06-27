package permissions

import (
	"context"
	"time"

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

func (s *svc) createPermission(ctx context.Context, payload CreatePermissionPayload) (PermissionResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return PermissionResponse{}, err
	}

	permission, err := s.repo.CreatePermission(ctx, repo.CreatePermissionParams{
		ID:           pgtype.UUID{Bytes: uid, Valid: true},
		StudentID:    pgtype.UUID{Bytes: uuid.MustParse(payload.StudentID), Valid: true},
		Type:         repo.PermissionType(payload.Type),
		Description:  payload.Description,
		RequestedAt:  pgtype.Timestamptz{Time: time.UnixMilli(payload.RequestedAt), Valid: true},
		ScheduledFor: pgtype.Timestamptz{Time: time.UnixMilli(payload.ScheduledFor), Valid: true},
	})
	if err != nil {
		return PermissionResponse{}, err
	}

	return permissionToResponse(permission), nil
}

func (s *svc) getPermissions(ctx context.Context, offset, limit int32) ([]PermissionResponse, error) {
	permissions, err := s.repo.GetAllPermissions(ctx, repo.GetAllPermissionsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return permissionsToResponse(permissions), nil
}

func (s *svc) getPermissionByID(ctx context.Context, id string) (PermissionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return PermissionResponse{}, err
	}
	permission, err := s.repo.GetPermissionByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return PermissionResponse{}, err
	}

	return permissionToResponse(permission), nil
}

func (s *svc) updatePermission(ctx context.Context, id string, payload UpdatePermissionPayload) (PermissionResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return PermissionResponse{}, err
	}

	// Read-modify-write: all fields are optional. Parse studentId only when
	// provided (uuid.MustParse would panic on empty) and keep existing
	// timestamps unless a new value was sent.
	existing, err := s.repo.GetPermissionByID(ctx, pgtype.UUID{Bytes: uidParsed, Valid: true})
	if err != nil {
		return PermissionResponse{}, err
	}

	params := repo.UpdatePermissionParams{
		ID:           existing.ID,
		StudentID:    existing.StudentID,
		Type:         existing.Type,
		Description:  existing.Description,
		RequestedAt:  existing.RequestedAt,
		ScheduledFor: existing.ScheduledFor,
	}
	if payload.StudentID != "" {
		studentID, err := uuid.Parse(payload.StudentID)
		if err != nil {
			return PermissionResponse{}, err
		}
		params.StudentID = pgtype.UUID{Bytes: studentID, Valid: true}
	}
	if payload.Type != "" {
		params.Type = repo.PermissionType(payload.Type)
	}
	if payload.Description != "" {
		params.Description = payload.Description
	}
	if payload.RequestedAt != 0 {
		params.RequestedAt = pgtype.Timestamptz{Time: time.UnixMilli(payload.RequestedAt), Valid: true}
	}
	if payload.ScheduledFor != 0 {
		params.ScheduledFor = pgtype.Timestamptz{Time: time.UnixMilli(payload.ScheduledFor), Valid: true}
	}

	permission, err := s.repo.UpdatePermission(ctx, params)
	if err != nil {
		return PermissionResponse{}, err
	}

	return permissionToResponse(permission), nil
}

func (s *svc) deletePermission(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeletePermission(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func permissionToResponse(permission repo.Permission) PermissionResponse {
	return PermissionResponse{
		ID:           permission.ID.String(),
		StudentID:    permission.StudentID.String(),
		Type:         string(permission.Type),
		Description:  string(permission.Description),
		RequestedAt:  permission.RequestedAt.Time.Format(time.RFC3339),
		ScheduledFor: permission.ScheduledFor.Time.Format(time.RFC3339),
	}
}

func permissionsToResponse(permissions []repo.Permission) []PermissionResponse {
	var response []PermissionResponse = make([]PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		response = append(response, permissionToResponse(permission))
	}
	return response
}

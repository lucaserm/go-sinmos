package studentguardians

import "context"

type (
	AddStudentGuardianPayload struct {
		StudentID  string `json:"studentId" validate:"required,uuid"`
		GuardianID string `json:"guardianId" validate:"required,uuid"`
	}

	RemoveStudentGuardianPayload struct {
		StudentID  string `json:"studentId" validate:"required,uuid"`
		GuardianID string `json:"guardianId" validate:"required,uuid"`
	}
)

type Service interface {
	addStudentGuardian(ctx context.Context, studentID, guardianID string) error
	removeStudentGuardian(ctx context.Context, studentID, guardianID string) error
}

package studentsubjects

import "context"

type (
	AddStudentSubjectPayload struct {
		StudentID string `json:"studentId" validate:"required,uuid"`
		SubjectID string `json:"subjectId" validate:"required,uuid"`
	}

	RemoveStudentSubjectPayload struct {
		StudentID string `json:"studentId" validate:"required,uuid"`
		SubjectID string `json:"subjectId" validate:"required,uuid"`
	}
)

type Service interface {
	addStudentSubject(ctx context.Context, studentID, subjectID string) error
	removeStudentSubject(ctx context.Context, studentID, subjectID string) error
}

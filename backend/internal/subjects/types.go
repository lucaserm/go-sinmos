package subjects

import "context"

type (
	CreateSubjectPayload struct {
		CourseID string `json:"courseId" validate:"required,uuid"`
		Name     string `json:"name" validate:"required"`
		Semester int32  `json:"semester" validate:"required"`
		Section  string `json:"section" validate:"required"`
	}

	UpdateSubjectPayload struct {
		CourseID string `json:"courseId,omitempty" validate:"omitempty,uuid"`
		Name     string `json:"name,omitempty"`
		Semester int32  `json:"semester,omitempty"`
		Section  string `json:"section,omitempty"`
	}

	SubjectResponse struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Semester int32  `json:"semester"`
		Section  string `json:"section"`
		CourseID string `json:"courseId"`
	}
)

type Service interface {
	createSubject(ctx context.Context, payload CreateSubjectPayload) (SubjectResponse, error)
	getSubjects(ctx context.Context, offset, limit int32) ([]SubjectResponse, error)
	getSubjectByID(ctx context.Context, id string) (SubjectResponse, error)
	updateSubject(ctx context.Context, id string, payload UpdateSubjectPayload) (SubjectResponse, error)
	deleteSubject(ctx context.Context, id string) error
}

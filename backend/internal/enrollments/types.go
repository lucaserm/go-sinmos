package enrollments

import "context"

type (
	CreateEnrollmentPayload struct {
		StudentID string `json:"studentId" validate:"required,uuid"`
		CourseID  string `json:"courseId" validate:"required,uuid"`
		Year      int32  `json:"year" validate:"required"`
	}

	UpdateEnrollmentPayload struct {
		StudentID string `json:"studentId,omitempty" validate:"omitempty,uuid"`
		CourseID  string `json:"courseId,omitempty" validate:"omitempty,uuid"`
		Year      int32  `json:"year,omitempty"`
	}

	EnrollmentResponse struct {
		ID        string `json:"id"`
		StudentID string `json:"studentId"`
		CourseID  string `json:"courseId"`
		Year      int32  `json:"year"`
	}
)

type Service interface {
	createEnrollment(ctx context.Context, payload CreateEnrollmentPayload) (EnrollmentResponse, error)
	getEnrollments(ctx context.Context, offset, limit int32) ([]EnrollmentResponse, error)
	getEnrollmentByID(ctx context.Context, id string) (EnrollmentResponse, error)
	updateEnrollment(ctx context.Context, id string, payload UpdateEnrollmentPayload) (EnrollmentResponse, error)
	deleteEnrollment(ctx context.Context, id string) error
}

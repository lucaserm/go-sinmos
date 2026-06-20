package courses

import "context"

type (
	CreateCoursePayload struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description,omitempty"`
		Session     string `json:"session" validate:"required,oneof=MORNING AFTERNOON EVENING"`
	}

	UpdateCoursePayload struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Session     string `json:"session,omitempty" validate:"oneof=MORNING AFTERNOON EVENING"`
	}

	CourseResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Session     string `json:"session"`
	}
)

type Service interface {
	createCourse(ctx context.Context, payload CreateCoursePayload) (CourseResponse, error)
	getCourses(ctx context.Context, offset, limit int32) ([]CourseResponse, error)
	getCourseByID(ctx context.Context, id string) (CourseResponse, error)
	updateCourse(ctx context.Context, id string, payload UpdateCoursePayload) (CourseResponse, error)
	deleteCourse(ctx context.Context, id string) error
}

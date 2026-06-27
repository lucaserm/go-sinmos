package courses

import (
	"context"

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

func (s *svc) createCourse(ctx context.Context, payload CreateCoursePayload) (CourseResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return CourseResponse{}, err
	}

	var description pgtype.Text

	if payload.Description != "" {
		description = pgtype.Text{String: payload.Description, Valid: true}
	} else {
		description = pgtype.Text{Valid: false}
	}

	course, err := s.repo.CreateCourse(ctx, repo.CreateCourseParams{
		ID:          pgtype.UUID{Bytes: uid, Valid: true},
		Name:        payload.Name,
		Description: description,
		Session:     repo.Session(payload.Session),
	})
	if err != nil {
		return CourseResponse{}, err
	}

	return courseToResponse(course), nil
}

func (s *svc) getCourses(ctx context.Context, offset, limit int32) ([]CourseResponse, error) {
	courses, err := s.repo.GetAllCourses(ctx, repo.GetAllCoursesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return coursesToResponse(courses), nil
}

func (s *svc) getCourseByID(ctx context.Context, id string) (CourseResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return CourseResponse{}, err
	}
	course, err := s.repo.GetCourseByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return CourseResponse{}, err
	}

	return courseToResponse(course), nil
}

func (s *svc) updateCourse(ctx context.Context, id string, payload UpdateCoursePayload) (CourseResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return CourseResponse{}, err
	}

	// Read-modify-write: every payload field is optional, so start from the
	// existing row and override only what was provided (avoids blanking fields).
	existing, err := s.repo.GetCourseByID(ctx, pgtype.UUID{Bytes: uidParsed, Valid: true})
	if err != nil {
		return CourseResponse{}, err
	}

	params := repo.UpdateCourseParams{
		ID:          existing.ID,
		Name:        existing.Name,
		Description: existing.Description,
		Session:     existing.Session,
	}
	if payload.Name != "" {
		params.Name = payload.Name
	}
	if payload.Description != "" {
		params.Description = pgtype.Text{String: payload.Description, Valid: true}
	}
	if payload.Session != "" {
		params.Session = repo.Session(payload.Session)
	}

	course, err := s.repo.UpdateCourse(ctx, params)
	if err != nil {
		return CourseResponse{}, err
	}

	return courseToResponse(course), nil
}

func (s *svc) deleteCourse(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteCourse(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func courseToResponse(course repo.Course) CourseResponse {
	return CourseResponse{
		ID:          course.ID.String(),
		Name:        course.Name,
		Description: course.Description.String,
		Session:     string(course.Session),
	}
}

func coursesToResponse(courses []repo.Course) []CourseResponse {
	var response []CourseResponse = make([]CourseResponse, 0, len(courses))
	for _, course := range courses {
		response = append(response, courseToResponse(course))
	}
	return response
}

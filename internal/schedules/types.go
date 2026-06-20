package schedules

import (
	"context"
)

/*

CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    session session NOT NULL,
    day_of_week VARCHAR(30) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
);*/

type (
	CreateSchedulePayload struct {
		SubjectID string `json:"subjectId" validate:"required"`
		Session   string `json:"session" validate:"required,oneof=MORNING AFTERNOON EVENING"`
		DayOfWeek string `json:"dayOfWeek" validate:"required"`
		StartTime string `json:"startTime" validate:"required"` // "08:00:00"
		EndTime   string `json:"endTime" validate:"required"`   // "10:00:00"
	}

	UpdateSchedulePayload struct {
		SubjectID string `json:"subjectId,omitempty"`
		Session   string `json:"session,omitempty" validate:"oneof=MORNING AFTERNOON EVENING"`
		DayOfWeek string `json:"dayOfWeek,omitempty"`
		StartTime string `json:"startTime,omitempty"` // "08:00:00"
		EndTime   string `json:"endTime,omitempty"`   // "10:00:00"
	}

	ScheduleResponse struct {
		ID        string `json:"id"`
		SubjectID string `json:"subjectId"`
		Session   string `json:"session"`
		DayOfWeek string `json:"dayOfWeek"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	}
)

type Service interface {
	createSchedule(ctx context.Context, payload CreateSchedulePayload) (ScheduleResponse, error)
	getSchedules(ctx context.Context, offset, limit int32) ([]ScheduleResponse, error)
	getScheduleByID(ctx context.Context, id string) (ScheduleResponse, error)
	updateSchedule(ctx context.Context, id string, payload UpdateSchedulePayload) (ScheduleResponse, error)
	deleteSchedule(ctx context.Context, id string) error
}

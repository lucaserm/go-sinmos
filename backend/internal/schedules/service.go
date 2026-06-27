package schedules

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/lucaserm/go-sinmos/internal/adapters/postgresql/sqlc"
	"github.com/lucaserm/go-sinmos/internal/utils"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) createSchedule(ctx context.Context, payload CreateSchedulePayload) (ScheduleResponse, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return ScheduleResponse{}, err
	}

	startParsed, err := time.Parse("15:04", payload.StartTime)
	if err != nil {
		return ScheduleResponse{}, err
	}

	endParsed, err := time.Parse("15:04", payload.EndTime)
	if err != nil {
		return ScheduleResponse{}, err
	}

	start := pgtype.Time{
		Microseconds: toMicroseconds(startParsed),
		Valid:        true,
	}

	end := pgtype.Time{
		Microseconds: toMicroseconds(endParsed),
		Valid:        true,
	}

	schedule, err := s.repo.CreateSchedule(ctx, repo.CreateScheduleParams{
		ID:        pgtype.UUID{Bytes: uid, Valid: true},
		SubjectID: pgtype.UUID{Bytes: uuid.MustParse(payload.SubjectID), Valid: true},
		Session:   repo.Session(payload.Session),
		DayOfWeek: payload.DayOfWeek,
		StartTime: start,
		EndTime:   end,
	})
	if err != nil {
		return ScheduleResponse{}, err
	}

	return scheduleToResponse(schedule), nil
}

func (s *svc) getSchedules(ctx context.Context, offset, limit int32) ([]ScheduleResponse, error) {
	schedules, err := s.repo.GetAllSchedules(ctx, repo.GetAllSchedulesParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	return schedulesToResponse(schedules), nil
}

func (s *svc) getScheduleByID(ctx context.Context, id string) (ScheduleResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ScheduleResponse{}, err
	}
	schedule, err := s.repo.GetScheduleByID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return ScheduleResponse{}, err
	}

	return scheduleToResponse(schedule), nil
}

func (s *svc) updateSchedule(ctx context.Context, id string, payload UpdateSchedulePayload) (ScheduleResponse, error) {
	uidParsed, err := uuid.Parse(id)
	if err != nil {
		return ScheduleResponse{}, err
	}

	// Read-modify-write: every field is optional. Parse subjectId/times only
	// when provided (uuid.MustParse and time.Parse would otherwise fail on empty
	// values and make a partial update impossible).
	existing, err := s.repo.GetScheduleByID(ctx, pgtype.UUID{Bytes: uidParsed, Valid: true})
	if err != nil {
		return ScheduleResponse{}, err
	}

	params := repo.UpdateScheduleParams{
		ID:        existing.ID,
		SubjectID: existing.SubjectID,
		Session:   existing.Session,
		DayOfWeek: existing.DayOfWeek,
		StartTime: existing.StartTime,
		EndTime:   existing.EndTime,
	}
	if payload.SubjectID != "" {
		subjectID, err := uuid.Parse(payload.SubjectID)
		if err != nil {
			return ScheduleResponse{}, err
		}
		params.SubjectID = pgtype.UUID{Bytes: subjectID, Valid: true}
	}
	if payload.Session != "" {
		params.Session = repo.Session(payload.Session)
	}
	if payload.DayOfWeek != "" {
		params.DayOfWeek = payload.DayOfWeek
	}
	if payload.StartTime != "" {
		startParsed, err := time.Parse("15:04", payload.StartTime)
		if err != nil {
			return ScheduleResponse{}, err
		}
		params.StartTime = pgtype.Time{Microseconds: toMicroseconds(startParsed), Valid: true}
	}
	if payload.EndTime != "" {
		endParsed, err := time.Parse("15:04", payload.EndTime)
		if err != nil {
			return ScheduleResponse{}, err
		}
		params.EndTime = pgtype.Time{Microseconds: toMicroseconds(endParsed), Valid: true}
	}

	schedule, err := s.repo.UpdateSchedule(ctx, params)
	if err != nil {
		return ScheduleResponse{}, err
	}

	return scheduleToResponse(schedule), nil
}

func (s *svc) deleteSchedule(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteSchedule(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

func toMicroseconds(t time.Time) int64 {
	return int64(
		t.Hour()*3600+
			t.Minute()*60+
			t.Second(),
	) * 1_000_000
}

func scheduleToResponse(schedule repo.Schedule) ScheduleResponse {
	return ScheduleResponse{
		ID:        schedule.ID.String(),
		SubjectID: schedule.SubjectID.String(),
		Session:   string(schedule.Session),
		DayOfWeek: schedule.DayOfWeek,
		StartTime: utils.PgTimeToString(schedule.StartTime),
		EndTime:   utils.PgTimeToString(schedule.EndTime),
	}
}

func schedulesToResponse(schedules []repo.Schedule) []ScheduleResponse {
	var response []ScheduleResponse = make([]ScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		response = append(response, scheduleToResponse(schedule))
	}
	return response
}

package engine

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// ScheduleManager реализует бизнес-логику для работы с расписанием
type ScheduleManager struct {
	*BaseManager[int, *domain.Schedule]
	db *sql.DB
}

func NewScheduleManager(
	repo db.Repository[*domain.Schedule, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ScheduleManager {
	return &ScheduleManager{
		BaseManager: NewBaseManager[int, *domain.Schedule](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByLesson возвращает расписание для конкретного занятия
func (m *ScheduleManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "lesson_id", Operator: "=", Value: lessonID},
		},
		OrderBy: "day_week, time_begin",
	})
	if err != nil {
		m.logger.Error("GetByLesson failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson_id", Value: lessonID},
		)
		return nil, fmt.Errorf("failed to get schedule by lesson: %w", err)
	}
	return DereferenceSlice(schedules), nil
}

// GetByDay возвращает расписание на конкретный день недели
func (m *ScheduleManager) GetByDay(ctx context.Context, dayWeek string) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "day_week", Operator: "=", Value: dayWeek},
		},
		OrderBy: "time_begin",
	})
	if err != nil {
		m.logger.Error("GetByDay failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "day_week", Value: dayWeek},
		)
		return nil, fmt.Errorf("failed to get schedule by day: %w", err)
	}
	return DereferenceSlice(schedules), nil
}

// GetCurrentSchedule возвращает актуальное расписание (в текущем периоде)
func (m *ScheduleManager) GetCurrentSchedule(ctx context.Context) ([]*domain.Schedule, error) {
	now := time.Now().Format("2006-01-02")
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "sched_date_start", Operator: "<=", Value: now},
			{Field: "sched_date_end", Operator: ">=", Value: now},
		},
		OrderBy: "day_week, time_begin",
	})
	if err != nil {
		m.logger.Error("GetCurrentSchedule failed",
			logger.Field{Key: "error", Value: err},
		)
		return nil, fmt.Errorf("failed to get current schedule: %w", err)
	}
	return DereferenceSlice(schedules), nil
}

// CheckTimeConflict проверяет наличие конфликтов в расписании
func (m *ScheduleManager) CheckTimeConflict(ctx context.Context, dayWeek, timeBegin, timeEnd string, excludeID int) (bool, error) {
	conflicts, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "day_week", Operator: "=", Value: dayWeek},
			{Field: "time_begin", Operator: "<", Value: timeEnd},
			{Field: "time_end", Operator: ">", Value: timeBegin},
			{Field: "schedule_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckTimeConflict failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "day_week", Value: dayWeek},
			logger.Field{Key: "time_begin", Value: timeBegin},
			logger.Field{Key: "time_end", Value: timeEnd},
		)
		return false, fmt.Errorf("failed to check schedule conflict: %w", err)
	}
	return len(conflicts) > 0, nil
}

// GetByDateRange возвращает расписание в указанном временном периоде
func (m *ScheduleManager) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "sched_date_start", Operator: "<=", Value: endDate.Format("2006-01-02")},
			{Field: "sched_date_end", Operator: ">=", Value: startDate.Format("2006-01-02")},
		},
		OrderBy: "sched_date_start, day_week, time_begin",
	})
	if err != nil {
		m.logger.Error("GetByDateRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "start_date", Value: startDate},
			logger.Field{Key: "end_date", Value: endDate},
		)
		return nil, fmt.Errorf("failed to get schedule by date range: %w", err)
	}
	return DereferenceSlice(schedules), nil
}

// GenerateSchedule генерирует расписание на основе шаблона (повторяющихся событий)
func (m *ScheduleManager) GenerateSchedule(ctx context.Context, template *domain.Schedule, until time.Time) error {
	currentDate := template.SchedDateStart
	dayMap := map[string]time.Weekday{
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
		"Sunday":    time.Sunday,
	}

	targetWeekday, ok := dayMap[template.DayWeek]
	if !ok {
		return fmt.Errorf("invalid day week: %s", template.DayWeek)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	for currentDate.Before(until) {
		for currentDate.Weekday() != targetWeekday {
			currentDate = currentDate.Add(24 * time.Hour)
		}

		if currentDate.After(until) {
			break
		}

		// Создаем новую запись расписания
		newSchedule := &domain.Schedule{
			LessonID:       template.LessonID,
			DayWeek:        template.DayWeek,
			TimeBegin:      template.TimeBegin,
			TimeEnd:        template.TimeEnd,
			SchedDateStart: currentDate,
			SchedDateEnd:   currentDate,
		}

		// Проверяем конфликты времени
		hasConflict, err := m.CheckTimeConflict(ctx, template.DayWeek, template.TimeBegin, template.TimeEnd, 0)
		if err != nil {
			return fmt.Errorf("failed to check time conflict: %w", err)
		}
		if hasConflict {
			return fmt.Errorf("time conflict detected for %s at %s-%s",
				template.DayWeek, template.TimeBegin, template.TimeEnd)
		}

		if err := txRepo.Create(ctx, &newSchedule); err != nil {
			return fmt.Errorf("failed to create schedule for date %s: %w",
				currentDate.Format("2006-01-02"), err)
		}

		currentDate = currentDate.Add(7 * 24 * time.Hour)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Create создает новую запись расписания
func (m *ScheduleManager) Create(ctx context.Context, schedule *domain.Schedule) error {
	if err := schedule.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	hasConflict, err := m.CheckTimeConflict(ctx, schedule.DayWeek, schedule.TimeBegin, schedule.TimeEnd, 0)
	if err != nil {
		return fmt.Errorf("failed to check time conflict: %w", err)
	}
	if hasConflict {
		return fmt.Errorf("time conflict detected for %s at %s-%s", schedule.DayWeek, schedule.TimeBegin, schedule.TimeEnd)
	}

	if err := m.repo.Create(ctx, &schedule); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// Update обновляет запись расписания
func (m *ScheduleManager) Update(ctx context.Context, schedule *domain.Schedule) error {
	if err := schedule.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	hasConflict, err := m.CheckTimeConflict(ctx, schedule.DayWeek, schedule.TimeBegin, schedule.TimeEnd, schedule.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to check time conflict: %w", err)
	}
	if hasConflict {
		return fmt.Errorf("time conflict detected for %s at %s-%s", schedule.DayWeek, schedule.TimeBegin, schedule.TimeEnd)
	}

	if err := m.repo.Update(ctx, &schedule); err != nil {
		m.logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

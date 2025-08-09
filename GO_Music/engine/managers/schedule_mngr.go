package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

// ScheduleManager реализует бизнес-логику для работы с расписанием
type ScheduleManager struct {
	*engine.BaseManager[int, domain.Schedule, *domain.Schedule]
	repo *repositories.ScheduleRepository
	db   *sql.DB
}

func NewScheduleManager(
	repo *repositories.ScheduleRepository,
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ScheduleManager {
	return &ScheduleManager{
		BaseManager: engine.NewBaseManager[int, domain.Schedule, *domain.Schedule](repo, logger, txTimeout),
		repo:        repo,
		db:          db,
	}
}

// [RU] GetByLesson возвращает расписание для конкретного занятия <--->
// [ENG] GetByLesson returns the schedule for a specific lesson
func (m *ScheduleManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "lesson_id", Operator: "=", Value: lessonID},
		},
		OrderBy: "day_week, time_begin",
	})
	if err != nil {
		m.Logger.Error("GetByLesson failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson_id", Value: lessonID},
		)
		return nil, fmt.Errorf("failed to get schedule by lesson: %w", err)
	}
	return schedules, nil
}

// [RU] GetByDay возвращает расписание на конкретный день недели <--->
// [ENG] GetByDay returns the schedule for a specific day of the week
func (m *ScheduleManager) GetByDay(ctx context.Context, dayWeek string) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "day_week", Operator: "=", Value: dayWeek},
		},
		OrderBy: "time_begin",
	})
	if err != nil {
		m.Logger.Error("GetByDay failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "day_week", Value: dayWeek},
		)
		return nil, fmt.Errorf("failed to get schedule by day: %w", err)
	}
	return schedules, nil
}

// [RU] GetCurrentSchedule возвращает актуальное расписание (в текущем периоде) <--->
// [ENG] GetCurrentSchedule returns the current schedule (in the current period)
func (m *ScheduleManager) GetCurrentSchedule(ctx context.Context) ([]*domain.Schedule, error) {
	now := time.Now().Format("2006-01-02")
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "schd_date_start", Operator: "<=", Value: now},
			{Field: "schd_date_end", Operator: ">=", Value: now},
		},
		OrderBy: "day_week, time_begin",
	})
	if err != nil {
		m.Logger.Error("GetCurrentSchedule failed",
			logger.Field{Key: "error", Value: err},
		)
		return nil, fmt.Errorf("failed to get current schedule: %w", err)
	}
	return schedules, nil
}

// [RU] CheckTimeConflict проверяет наличие конфликтов в расписании <--->
// [ENG] CheckTimeConflict checks for conflicts in the schedule
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
		m.Logger.Error("CheckTimeConflict failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "day_week", Value: dayWeek},
			logger.Field{Key: "time_begin", Value: timeBegin},
			logger.Field{Key: "time_end", Value: timeEnd},
		)
		return false, fmt.Errorf("failed to check schedule conflict: %w", err)
	}
	return len(conflicts) > 0, nil
}

// [RU] GetByDateRange возвращает расписание в указанном временном периоде <--->
// [ENG] GetByDateRange returns the schedule in the specified date range
func (m *ScheduleManager) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "schd_date_start", Operator: "<=", Value: endDate.Format("2006-01-02")},
			{Field: "schd_date_end", Operator: ">=", Value: startDate.Format("2006-01-02")},
		},
		OrderBy: "schd_date_start, day_week, time_begin",
	})
	if err != nil {
		m.Logger.Error("GetByDateRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "start_date", Value: startDate},
			logger.Field{Key: "end_date", Value: endDate},
		)
		return nil, fmt.Errorf("failed to get schedule by date range: %w", err)
	}
	return schedules, nil
}

// [RU] GenerateSchedule генерирует расписание на основе шаблона (повторяющихся событий) <--->
// [ENG] GenerateSchedule generates a schedule based on a template (recurring events)
func (m *ScheduleManager) GenerateSchedule(ctx context.Context, template *domain.Schedule, until time.Time) error {
	currentDate := template.SchdDateStart
	dayMap := map[string]time.Weekday{
		"Понедельник": time.Monday,
		"Вторник":     time.Tuesday,
		"Среда":       time.Wednesday,
		"Четверг":     time.Thursday,
		"Пятница":     time.Friday,
		"Суббота":     time.Saturday,
		"Воскресенье": time.Sunday,
	}

	targetWeekday, ok := dayMap[template.DayWeek]
	if !ok {
		return fmt.Errorf("invalid day week: %s", template.DayWeek)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.Repo.WithTx(tx)

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
			LessonID:      template.LessonID,
			DayWeek:       template.DayWeek,
			TimeBegin:     template.TimeBegin,
			TimeEnd:       template.TimeEnd,
			SchdDateStart: currentDate,
			SchdDateEnd:   currentDate,
		}

		// Проверяем конфликты времени
		hasConflict, err := m.CheckTimeConflict(ctx, template.DayWeek, template.TimeBegin.Format("15:04"), template.TimeEnd.Format("15:04"), 0)
		if err != nil {
			return fmt.Errorf("failed to check time conflict: %w", err)
		}
		if hasConflict {
			return fmt.Errorf("time conflict detected for %s at %s-%s",
				template.DayWeek, template.TimeBegin.Format("15:04"), template.TimeEnd.Format("15:04"))
		}

		if err := txRepo.Create(ctx, newSchedule); err != nil {
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

// [RU] Create создает новую запись расписания <--->
// [ENG] Create creates a new schedule entry
func (m *ScheduleManager) Create(ctx context.Context, schedule *domain.Schedule) error {
	if err := schedule.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	hasConflict, err := m.CheckTimeConflict(ctx, schedule.DayWeek, schedule.TimeBegin.Format("15:04"), schedule.TimeEnd.Format("15:04"), 0)
	if err != nil {
		return fmt.Errorf("failed to check time conflict: %w", err)
	}
	if hasConflict {
		return fmt.Errorf("time conflict detected for %s at %s-%s", schedule.DayWeek, schedule.TimeBegin.Format("15:04"), schedule.TimeEnd.Format("15:04"))
	}

	if err := m.Repo.Create(ctx, schedule); err != nil {
		m.Logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// [RU] Update обновляет запись расписания <--->
// [ENG] Update updates the schedule entry
func (m *ScheduleManager) Update(ctx context.Context, schedule *domain.Schedule) error {
	if err := schedule.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	hasConflict, err := m.CheckTimeConflict(ctx, schedule.DayWeek, schedule.TimeBegin.Format("15:04"), schedule.TimeEnd.Format("15:04"), schedule.ScheduleID)
	if err != nil {
		return fmt.Errorf("failed to check time conflict: %w", err)
	}
	if hasConflict {
		return fmt.Errorf("time conflict detected for %s at %s-%s", schedule.DayWeek, schedule.TimeBegin.Format("15:04"), schedule.TimeEnd.Format("15:04"))
	}

	if err := m.Repo.Update(ctx, schedule); err != nil {
		m.Logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "schedule", Value: schedule},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

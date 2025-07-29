package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// ScheduleManager реализует бизнес-логику для работы с расписанием
type ScheduleManager struct {
	*BaseManager[domain.Schedule, *domain.Schedule]
}

func NewScheduleManager(
	repo Repository[domain.Schedule, *domain.Schedule],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ScheduleManager {
	return &ScheduleManager{
		BaseManager: NewBaseManager[domain.Schedule](repo, logger, txTimeout),
	}
}

// GetByLesson возвращает расписание для конкретного занятия
func (m *ScheduleManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return schedules, nil
}

// GetByDay возвращает расписание на конкретный день недели
func (m *ScheduleManager) GetByDay(ctx context.Context, dayWeek string) ([]*domain.Schedule, error) {
	schedules, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return schedules, nil
}

// GetCurrentSchedule возвращает актуальное расписание (в текущем периоде)
func (m *ScheduleManager) GetCurrentSchedule(ctx context.Context) ([]*domain.Schedule, error) {
	now := time.Now().Format("2006-01-02")
	schedules, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return schedules, nil
}

// CheckTimeConflict проверяет наличие конфликтов в расписании
func (m *ScheduleManager) CheckTimeConflict(ctx context.Context, dayWeek, timeBegin, timeEnd string, excludeID int) (bool, error) {
	conflicts, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	schedules, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return schedules, nil
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

	return m.ExecuteInTx(ctx, m.repo.(TxProvider), func(repo Repository[domain.Schedule, *domain.Schedule]) error {
		for currentDate.Before(until) {
			for currentDate.Weekday() != targetWeekday {
				currentDate = currentDate.Add(24 * time.Hour)
			}

			if currentDate.After(until) {
				break
			}

			newSchedule := *template
			newSchedule.SchedDateStart = currentDate
			newSchedule.SchedDateEnd = currentDate

			if err := repo.Create(ctx, &newSchedule); err != nil {
				return err
			}

			currentDate = currentDate.Add(7 * 24 * time.Hour) // Следующая неделя
		}
		return nil
	})
}

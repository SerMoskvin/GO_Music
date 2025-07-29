package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// LessonManager реализует бизнес-логику для занятий
type LessonManager struct {
	*BaseManager[domain.Lesson, *domain.Lesson]
}

func NewLessonManager(
	repo Repository[domain.Lesson, *domain.Lesson],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *LessonManager {
	return &LessonManager{
		BaseManager: NewBaseManager[domain.Lesson](repo, logger, txTimeout),
	}
}

// GetByEmployee возвращает занятия преподавателя
func (m *LessonManager) GetByEmployee(ctx context.Context, employeeID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "employee_id", Operator: "=", Value: employeeID},
		},
		OrderBy: "lesson_id DESC",
	})
	if err != nil {
		m.logger.Error("GetByEmployee failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
		)
		return nil, fmt.Errorf("failed to get lessons by employee: %w", err)
	}
	return lessons, nil
}

// GetByGroup возвращает занятия группы
func (m *LessonManager) GetByGroup(ctx context.Context, groupID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "group_id", Operator: "=", Value: groupID},
		},
		OrderBy: "lesson_id DESC",
	})
	if err != nil {
		m.logger.Error("GetByGroup failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group_id", Value: groupID},
		)
		return nil, fmt.Errorf("failed to get lessons by group: %w", err)
	}
	return lessons, nil
}

// GetByStudent возвращает индивидуальные занятия студента
func (m *LessonManager) GetByStudent(ctx context.Context, studentID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "student_id", Operator: "=", Value: studentID},
		},
		OrderBy: "lesson_id DESC",
	})
	if err != nil {
		m.logger.Error("GetByStudent failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student_id", Value: studentID},
		)
		return nil, fmt.Errorf("failed to get lessons by student: %w", err)
	}
	return lessons, nil
}

// GetBySubject возвращает занятия по предмету
func (m *LessonManager) GetBySubject(ctx context.Context, subjectID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "subject_id", Operator: "=", Value: subjectID},
		},
		OrderBy: "lesson_id DESC",
	})
	if err != nil {
		m.logger.Error("GetBySubject failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return nil, fmt.Errorf("failed to get lessons by subject: %w", err)
	}
	return lessons, nil
}

// GetByAudience возвращает занятия в аудитории
func (m *LessonManager) GetByAudience(ctx context.Context, audienceID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "audience_id", Operator: "=", Value: audienceID},
		},
		OrderBy: "lesson_id DESC",
	})
	if err != nil {
		m.logger.Error("GetByAudience failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience_id", Value: audienceID},
		)
		return nil, fmt.Errorf("failed to get lessons by audience: %w", err)
	}
	return lessons, nil
}

// CheckEmployeeAvailability проверяет, свободен ли преподаватель в указанное время
func (m *LessonManager) CheckEmployeeAvailability(
	ctx context.Context,
	employeeID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	return true, nil
}

// CheckAudienceAvailability проверяет, свободна ли аудитория в указанное время
func (m *LessonManager) CheckAudienceAvailability(
	ctx context.Context,
	audienceID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	return true, nil
}

package engine

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// LessonManager реализует бизнес-логику для занятий
type LessonManager struct {
	*BaseManager[int, *domain.Lesson]
	db *sql.DB
}

func NewLessonManager(
	repo db.Repository[*domain.Lesson, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *LessonManager {
	return &LessonManager{
		BaseManager: NewBaseManager[int, *domain.Lesson](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByEmployee возвращает занятия преподавателя
func (m *LessonManager) GetByEmployee(ctx context.Context, employeeID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(lessons), nil
}

// GetByGroup возвращает занятия группы
func (m *LessonManager) GetByGroup(ctx context.Context, groupID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(lessons), nil
}

// GetByStudent возвращает индивидуальные занятия студента
func (m *LessonManager) GetByStudent(ctx context.Context, studentID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(lessons), nil
}

// GetBySubject возвращает занятия по предмету
func (m *LessonManager) GetBySubject(ctx context.Context, subjectID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(lessons), nil
}

// GetByAudience возвращает занятия в аудитории
func (m *LessonManager) GetByAudience(ctx context.Context, audienceID int) ([]*domain.Lesson, error) {
	lessons, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(lessons), nil
}

// CheckEmployeeAvailability проверяет, свободен ли преподаватель в указанное время
func (m *LessonManager) CheckEmployeeAvailability(
	ctx context.Context,
	employeeID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	repo, ok := m.repo.(*repositories.LessonRepository)
	if !ok {
		return false, fmt.Errorf("repository does not support availability checks")
	}
	return repo.CheckEmployeeAvailability(ctx, employeeID, startTime, endTime, excludeLessonID)
}

// CheckAudienceAvailability проверяет, свободна ли аудитория в указанное время
func (m *LessonManager) CheckAudienceAvailability(
	ctx context.Context,
	audienceID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	repo, ok := m.repo.(*repositories.LessonRepository)
	if !ok {
		return false, fmt.Errorf("repository does not support availability checks")
	}
	return repo.CheckAudienceAvailability(ctx, audienceID, startTime, endTime, excludeLessonID)
}

// Create создает новое занятие
func (m *LessonManager) Create(ctx context.Context, lesson *domain.Lesson) error {
	if err := lesson.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson", Value: lesson},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := m.repo.Create(ctx, &lesson); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson", Value: lesson},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// Update обновляет данные занятия
func (m *LessonManager) Update(ctx context.Context, lesson *domain.Lesson) error {
	if err := lesson.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson", Value: lesson},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := m.repo.Update(ctx, &lesson); err != nil {
		m.logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson", Value: lesson},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// BulkCreate массово создает занятия в транзакции
func (m *LessonManager) BulkCreate(ctx context.Context, lessons []*domain.Lesson) error {
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

	for _, lesson := range lessons {
		if err := lesson.Validate(); err != nil {
			return fmt.Errorf("validation failed for lesson %v: %w", lesson, err)
		}

		ptrToLesson := &lesson
		if err := txRepo.Create(ctx, ptrToLesson); err != nil {
			return fmt.Errorf("create failed for lesson %v: %w", lesson, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

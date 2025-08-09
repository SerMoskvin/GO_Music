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
	*BaseManager[int, domain.Lesson, *domain.Lesson]
	repo *repositories.LessonRepository
	db   *sql.DB
}

func NewLessonManager(
	repo *repositories.LessonRepository,
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *LessonManager {
	return &LessonManager{
		BaseManager: NewBaseManager[int, domain.Lesson, *domain.Lesson](repo, logger, txTimeout),
		repo:        repo,
		db:          db,
	}
}

// [RU] GetByEmployee возвращает занятия преподавателя <--->
// [ENG] GetByEmployee returns lessons for the specified employee
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
	return lessons, nil
}

// [RU] GetByGroup возвращает занятия группы <--->
// [ENG] GetByGroup returns lessons for the specified group
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
	return lessons, nil
}

// [RU] GetByStudent возвращает индивидуальные занятия студента <--->
// [ENG] GetByStudent returns individual lessons for the specified student
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
	return lessons, nil
}

// [RU] GetBySubject возвращает занятия по предмету <--->
// [ENG] GetBySubject returns lessons for the specified subject
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
	return lessons, nil
}

// [RU] GetByAudience возвращает занятия в аудитории <--->
// [ENG] GetByAudience returns lessons for the specified audience
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
	return lessons, nil
}

// [RU] CheckEmployeeAvailability проверяет, свободен ли преподаватель в указанное время <--->
// [ENG] CheckEmployeeAvailability checks if the employee is available during the specified time
func (m *LessonManager) CheckEmployeeAvailability(
	ctx context.Context,
	employeeID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	return m.repo.CheckEmployeeAvailability(ctx, employeeID, startTime, endTime, excludeLessonID)
}

// [RU] CheckAudienceAvailability проверяет, свободна ли аудитория в указанное время <--->
// [ENG] CheckAudienceAvailability checks if the audience is available during the specified time
func (m *LessonManager) CheckAudienceAvailability(
	ctx context.Context,
	audienceID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	return m.repo.CheckAudienceAvailability(ctx, audienceID, startTime, endTime, excludeLessonID)
}

// [RU] BulkCreate массово создает занятия в транзакции <--->
// [ENG] BulkCreate creates multiple lessons in a transaction
func (m *LessonManager) BulkCreate(ctx context.Context, lessons []*domain.Lesson) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, lesson := range lessons {
		if err := lesson.Validate(); err != nil {
			return fmt.Errorf("validation failed for lesson %v: %w", lesson, err)
		}

		if err := txRepo.Create(ctx, lesson); err != nil {
			return fmt.Errorf("create failed for lesson %v: %w", lesson, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

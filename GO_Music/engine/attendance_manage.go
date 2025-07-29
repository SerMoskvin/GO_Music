package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// StudentAttendanceManager реализует бизнес-логику для посещаемости студентов
type StudentAttendanceManager struct {
	*BaseManager[domain.StudentAttendance, *domain.StudentAttendance]
}

func NewStudentAttendanceManager(
	repo Repository[domain.StudentAttendance, *domain.StudentAttendance],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentAttendanceManager {
	return &StudentAttendanceManager{
		BaseManager: NewBaseManager[domain.StudentAttendance](repo, logger, txTimeout),
	}
}

// GetByStudent возвращает записи посещаемости для конкретного студента
func (m *StudentAttendanceManager) GetByStudent(ctx context.Context, studentID int) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "student_id", Operator: "=", Value: studentID},
		},
		OrderBy: "attendance_date DESC",
	})
	if err != nil {
		m.logger.Error("GetByStudent failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student_id", Value: studentID},
		)
		return nil, fmt.Errorf("failed to get attendance by student: %w", err)
	}
	return records, nil
}

// GetByLesson возвращает записи посещаемости для конкретного занятия
func (m *StudentAttendanceManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "lesson_id", Operator: "=", Value: lessonID},
		},
		OrderBy: "student_id",
	})
	if err != nil {
		m.logger.Error("GetByLesson failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson_id", Value: lessonID},
		)
		return nil, fmt.Errorf("failed to get attendance by lesson: %w", err)
	}
	return records, nil
}

// GetByDateRange возвращает записи за указанный период
func (m *StudentAttendanceManager) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "attendance_date", Operator: ">=", Value: startDate},
			{Field: "attendance_date", Operator: "<=", Value: endDate},
		},
		OrderBy: "attendance_date, student_id",
	})
	if err != nil {
		m.logger.Error("GetByDateRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "start_date", Value: startDate},
			logger.Field{Key: "end_date", Value: endDate},
		)
		return nil, fmt.Errorf("failed to get attendance by date range: %w", err)
	}
	return records, nil
}

// GetStudentAttendanceStats возвращает статистику посещаемости студента
func (m *StudentAttendanceManager) GetStudentAttendanceStats(ctx context.Context, studentID int) (present, absent int, err error) {
	records, err := m.GetByStudent(ctx, studentID)
	if err != nil {
		return 0, 0, err
	}

	for _, r := range records {
		if r.PresenceMark {
			present++
		} else {
			absent++
		}
	}
	return present, absent, nil
}

// BulkCreate создает несколько записей посещаемости в транзакции
func (m *StudentAttendanceManager) BulkCreate(ctx context.Context, records []*domain.StudentAttendance) error {
	return m.ExecuteInTx(ctx, m.repo.(TxProvider), func(repo Repository[domain.StudentAttendance, *domain.StudentAttendance]) error {
		for _, record := range records {
			if err := record.Validate(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
			if err := repo.Create(ctx, record); err != nil {
				return fmt.Errorf("create failed: %w", err)
			}
		}
		return nil
	})
}

// CheckDuplicate проверяет наличие дублирующей записи посещаемости
func (m *StudentAttendanceManager) CheckDuplicate(ctx context.Context, studentID, lessonID int) (bool, error) {
	records, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "student_id", Operator: "=", Value: studentID},
			{Field: "lesson_id", Operator: "=", Value: lessonID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckDuplicate failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student_id", Value: studentID},
			logger.Field{Key: "lesson_id", Value: lessonID},
		)
		return false, fmt.Errorf("failed to check attendance duplicate: %w", err)
	}
	return len(records) > 0, nil
}

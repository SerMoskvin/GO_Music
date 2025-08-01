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

// StudentAttendanceManager реализует бизнес-логику для посещаемости студентов
type StudentAttendanceManager struct {
	*BaseManager[int, *domain.StudentAttendance]
	db *sql.DB
}

func NewStudentAttendanceManager(
	repo db.Repository[*domain.StudentAttendance, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentAttendanceManager {
	return &StudentAttendanceManager{
		BaseManager: NewBaseManager[int, *domain.StudentAttendance](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByStudent возвращает записи посещаемости для конкретного студента
func (m *StudentAttendanceManager) GetByStudent(ctx context.Context, studentID int) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(records), nil
}

// GetByLesson возвращает записи посещаемости для конкретного занятия
func (m *StudentAttendanceManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(records), nil
}

// GetByDateRange возвращает записи за указанный период
func (m *StudentAttendanceManager) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(records), nil
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

	for _, record := range records {
		if err := record.Validate(); err != nil {
			return fmt.Errorf("validation failed for record %v: %w", record, err)
		}

		// Создаем указатель на указатель, который ожидает репозиторий
		ptrToRecord := &record
		if err := txRepo.Create(ctx, ptrToRecord); err != nil {
			return fmt.Errorf("create failed for record %v: %w", record, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CheckDuplicate проверяет наличие дублирующей записи посещаемости
func (m *StudentAttendanceManager) CheckDuplicate(ctx context.Context, studentID, lessonID int) (bool, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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

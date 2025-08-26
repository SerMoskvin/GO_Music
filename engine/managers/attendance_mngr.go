package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"
	e "GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

type StudentAttendanceManager struct {
	*e.BaseManager[int, domain.StudentAttendance, *domain.StudentAttendance]
	db *sql.DB
}

func NewStudentAttendanceManager(
	repo db.Repository[domain.StudentAttendance, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentAttendanceManager {
	return &StudentAttendanceManager{
		BaseManager: e.NewBaseManager[int, domain.StudentAttendance, *domain.StudentAttendance](repo, logger, txTimeout),
		db:          db,
	}
}

// [RU] GetByStudent возвращает записи посещаемости для конкретного студента <--->
// [ENG] GetByStudent returns all student's attendance records
func (m *StudentAttendanceManager) GetByStudent(ctx context.Context, studentID int) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "student_id", Operator: "=", Value: studentID},
		},
		OrderBy: "attendance_date DESC",
	})
	if err != nil {
		m.Logger.Error("GetByStudent failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student_id", Value: studentID},
		)
		return nil, fmt.Errorf("failed to get attendance by student: %w", err)
	}
	return records, nil
}

// [RU] GetByLesson возвращает записи посещаемости для конкретного занятия <--->
// [ENG] GetByLesson returns all attendance records for a specific lesson
func (m *StudentAttendanceManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "lesson_id", Operator: "=", Value: lessonID},
		},
		OrderBy: "student_id",
	})
	if err != nil {
		m.Logger.Error("GetByLesson failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "lesson_id", Value: lessonID},
		)
		return nil, fmt.Errorf("failed to get attendance by lesson: %w", err)
	}
	return records, nil
}

// [RU] GetByDateRange возвращает записи за указанный период <--->
// [ENG] GetByDateRange returns attendance records for the specified date range
func (m *StudentAttendanceManager) GetByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.StudentAttendance, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "attendance_date", Operator: ">=", Value: startDate},
			{Field: "attendance_date", Operator: "<=", Value: endDate},
		},
		OrderBy: "attendance_date, student_id",
	})
	if err != nil {
		m.Logger.Error("GetByDateRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "start_date", Value: startDate},
			logger.Field{Key: "end_date", Value: endDate},
		)
		return nil, fmt.Errorf("failed to get attendance by date range: %w", err)
	}
	return records, nil
}

// [RU] GetStudentAttendanceStats возвращает статистику посещаемости студента <--->
// [ENG] GetStudentAttendanceStats returns the attendance statistics for a student
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

// [RU] BulkCreate создает несколько записей посещаемости в транзакции <--->
// [ENG] BulkCreate creates multiple attendance records in a transaction
func (m *StudentAttendanceManager) BulkCreate(ctx context.Context, records []*domain.StudentAttendance) error {
	return m.ExecuteInTx(ctx, m.db, func(txRepo db.Repository[domain.StudentAttendance, int]) error {
		for _, record := range records {
			if err := record.Validate(); err != nil {
				return fmt.Errorf("validation failed for record %v: %w", record, err)
			}
			if err := txRepo.Create(ctx, record); err != nil {
				return fmt.Errorf("create failed for record %v: %w", record, err)
			}
		}
		return nil
	})
}

// [RU] CheckDuplicate проверяет наличие дублирующей записи посещаемости <--->
// [ENG] CheckDuplicate checks for the existence of a duplicate attendance record
func (m *StudentAttendanceManager) CheckDuplicate(ctx context.Context, studentID, lessonID int) (bool, error) {
	records, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "student_id", Operator: "=", Value: studentID},
			{Field: "lesson_id", Operator: "=", Value: lessonID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckDuplicate failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student_id", Value: studentID},
			logger.Field{Key: "lesson_id", Value: lessonID},
		)
		return false, fmt.Errorf("failed to check attendance duplicate: %w", err)
	}
	return len(records) > 0, nil
}

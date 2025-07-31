package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// StudentAssessmentManager реализует бизнес-логику для оценок студентов
type StudentAssessmentManager struct {
	*BaseManager[int, *domain.StudentAssessment]
}

func NewStudentAssessmentManager(
	repo db.Repository[*domain.StudentAssessment, int],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentAssessmentManager {
	return &StudentAssessmentManager{
		BaseManager: NewBaseManager[int, *domain.StudentAssessment](repo, logger, txTimeout),
	}
}

// GetByStudent возвращает все оценки студента
func (m *StudentAssessmentManager) GetByStudent(ctx context.Context, studentID int) ([]*domain.StudentAssessment, error) {
	assessments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "student_id", Operator: "=", Value: studentID},
		},
		OrderBy: "assessment_date DESC",
	})
	if err != nil {
		m.logger.Error("GetByStudent failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student_id", Value: studentID},
		)
		return nil, fmt.Errorf("failed to get assessments by student: %w", err)
	}
	return assessments, nil
}

// GetByLesson возвращает все оценки за конкретное занятие
func (m *StudentAssessmentManager) GetByLesson(ctx context.Context, lessonID int) ([]*domain.StudentAssessment, error) {
	assessments, err := m.List(ctx, Filter{
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
		return nil, fmt.Errorf("failed to get assessments by lesson: %w", err)
	}
	return assessments, nil
}

// GetByTaskType возвращает оценки по типу задания
func (m *StudentAssessmentManager) GetByTaskType(ctx context.Context, taskType string) ([]*domain.StudentAssessment, error) {
	assessments, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "task_type", Operator: "=", Value: taskType},
		},
		OrderBy: "assessment_date DESC",
	})
	if err != nil {
		m.logger.Error("GetByTaskType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "task_type", Value: taskType},
		)
		return nil, fmt.Errorf("failed to get assessments by task type: %w", err)
	}
	return assessments, nil
}

// GetStudentAverageGrade вычисляет средний балл студента
func (m *StudentAssessmentManager) GetStudentAverageGrade(ctx context.Context, studentID int) (float64, error) {
	assessments, err := m.GetByStudent(ctx, studentID)
	if err != nil {
		return 0, err
	}

	if len(assessments) == 0 {
		return 0, nil
	}

	var sum int
	for _, a := range assessments {
		sum += a.Grade
	}

	return float64(sum) / float64(len(assessments)), nil
}

// GetGradesByDateRange возвращает оценки за период
func (m *StudentAssessmentManager) GetGradesByDateRange(ctx context.Context, startDate, endDate string) ([]*domain.StudentAssessment, error) {
	assessments, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "assessment_date", Operator: ">=", Value: startDate},
			{Field: "assessment_date", Operator: "<=", Value: endDate},
		},
		OrderBy: "assessment_date, student_id",
	})
	if err != nil {
		m.logger.Error("GetGradesByDateRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "start_date", Value: startDate},
			logger.Field{Key: "end_date", Value: endDate},
		)
		return nil, fmt.Errorf("failed to get assessments by date range: %w", err)
	}
	return assessments, nil
}

// BulkUpsert массовое обновление/добавление оценок в транзакции
func (m *StudentAssessmentManager) BulkUpsert(ctx context.Context, assessments []*domain.StudentAssessment) error {
	return m.ExecuteInTx(ctx, m.repo.(TxProvider), func(repo Repository[domain.StudentAssessment, *domain.StudentAssessment]) error {
		for _, assessment := range assessments {
			if err := assessment.Validate(); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// Проверяем существование записи
			exists, err := repo.Exists(ctx, assessment.AssessmentNoteID)
			if err != nil {
				return fmt.Errorf("exists check failed: %w", err)
			}

			if exists {
				if err := repo.Update(ctx, assessment); err != nil {
					return fmt.Errorf("update failed: %w", err)
				}
			} else {
				if err := repo.Create(ctx, assessment); err != nil {
					return fmt.Errorf("create failed: %w", err)
				}
			}
		}
		return nil
	})
}

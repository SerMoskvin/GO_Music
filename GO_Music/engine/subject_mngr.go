package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// SubjectManager реализует бизнес-логику для работы с предметами
type SubjectManager struct {
	*BaseManager[domain.Subject, *domain.Subject]
}

func NewSubjectManager(
	repo Repository[domain.Subject, *domain.Subject],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *SubjectManager {
	return &SubjectManager{
		BaseManager: NewBaseManager[domain.Subject](repo, logger, txTimeout),
	}
}

// GetByType возвращает предметы указанного типа
func (m *SubjectManager) GetByType(ctx context.Context, subjectType string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "subject_type", Operator: "=", Value: subjectType},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.logger.Error("GetByType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "type", Value: subjectType},
		)
		return nil, fmt.Errorf("failed to get subjects by type: %w", err)
	}
	return subjects, nil
}

// SearchByName ищет предметы по названию
func (m *SubjectManager) SearchByName(ctx context.Context, name string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "subject_name", Operator: "ILIKE", Value: "%" + name + "%"},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.logger.Error("SearchByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return nil, fmt.Errorf("failed to search subjects by name: %w", err)
	}
	return subjects, nil
}

// GetByDescription ищет предметы по описанию
func (m *SubjectManager) GetByDescription(ctx context.Context, keyword string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "short_desc", Operator: "ILIKE", Value: "%" + keyword + "%"},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.logger.Error("GetByDescription failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "keyword", Value: keyword},
		)
		return nil, fmt.Errorf("failed to get subjects by description: %w", err)
	}
	return subjects, nil
}

// CheckNameUnique проверяет уникальность названия предмета
func (m *SubjectManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	subjects, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "subject_name", Operator: "=", Value: name},
			{Field: "subject_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(subjects) == 0, nil
}

// GetSubjectsWithPrograms возвращает предметы с привязанными программами
func (m *SubjectManager) GetSubjectsWithPrograms(ctx context.Context, programID int) ([]*domain.Subject, error) {
	// В реальной реализации здесь будет JOIN с таблицей связей предметов и программ
	subjects, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programID},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.logger.Error("GetSubjectsWithPrograms failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "program_id", Value: programID},
		)
		return nil, fmt.Errorf("failed to get subjects with programs: %w", err)
	}
	return subjects, nil
}

// GetPopularSubjects возвращает самые популярные предметы (по количеству студентов)
func (m *SubjectManager) GetPopularSubjects(ctx context.Context, limit int) ([]*domain.Subject, error) {
	// В реальной реализации здесь будет сложный запрос с подсчетом студентов
	subjects, err := m.List(ctx, Filter{
		OrderBy: "student_count DESC",
		Limit:   limit,
	})
	if err != nil {
		m.logger.Error("GetPopularSubjects failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "limit", Value: limit},
		)
		return nil, fmt.Errorf("failed to get popular subjects: %w", err)
	}
	return subjects, nil
}

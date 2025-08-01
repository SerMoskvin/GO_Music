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

// SubjectManager реализует бизнес-логику для работы с предметами
type SubjectManager struct {
	*BaseManager[int, *domain.Subject]
	db *sql.DB
}

func NewSubjectManager(
	repo db.Repository[*domain.Subject, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *SubjectManager {
	return &SubjectManager{
		BaseManager: NewBaseManager[int, *domain.Subject](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByType возвращает предметы указанного типа
func (m *SubjectManager) GetByType(ctx context.Context, subjectType string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(subjects), nil
}

// SearchByName ищет предметы по названию
func (m *SubjectManager) SearchByName(ctx context.Context, name string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(subjects), nil
}

// GetByDescription ищет предметы по описанию
func (m *SubjectManager) GetByDescription(ctx context.Context, keyword string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(subjects), nil
}

// CheckNameUnique проверяет уникальность названия предмета
func (m *SubjectManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	repo, ok := m.repo.(interface {
		GetSubjectsWithPrograms(ctx context.Context, programID int) ([]*domain.Subject, error)
	})
	if !ok {
		return nil, fmt.Errorf("repository doesn't support GetSubjectsWithPrograms")
	}

	subjects, err := repo.GetSubjectsWithPrograms(ctx, programID)
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
	repo, ok := m.repo.(interface {
		GetPopularSubjects(ctx context.Context, limit int) ([]*domain.Subject, error)
	})
	if !ok {
		return nil, fmt.Errorf("repository doesn't support GetPopularSubjects")
	}

	subjects, err := repo.GetPopularSubjects(ctx, limit)
	if err != nil {
		m.logger.Error("GetPopularSubjects failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "limit", Value: limit},
		)
		return nil, fmt.Errorf("failed to get popular subjects: %w", err)
	}
	return subjects, nil
}

// Create создает новый предмет
func (m *SubjectManager) Create(ctx context.Context, subject *domain.Subject) error {
	if err := subject.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject", Value: subject},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, subject.SubjectName, 0)
	if err != nil {
		return fmt.Errorf("uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("subject name %s already exists", subject.SubjectName)
	}

	ptrToSubject := &subject
	if err := m.repo.Create(ctx, ptrToSubject); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject", Value: subject},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// Update обновляет предмет
func (m *SubjectManager) Update(ctx context.Context, subject *domain.Subject) error {
	if err := subject.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject", Value: subject},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, subject.SubjectName, subject.SubjectID)
	if err != nil {
		return fmt.Errorf("uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("subject name %s already exists", subject.SubjectName)
	}

	ptrToSubject := &subject
	if err := m.repo.Update(ctx, ptrToSubject); err != nil {
		m.logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject", Value: subject},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

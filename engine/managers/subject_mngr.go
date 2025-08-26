package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

type SubjectManager struct {
	*engine.BaseManager[int, domain.Subject, *domain.Subject]
	repo *repositories.SubjectRepository
	db   *sql.DB
}

func NewSubjectManager(
	repo *repositories.SubjectRepository,
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *SubjectManager {
	return &SubjectManager{
		BaseManager: engine.NewBaseManager[int, domain.Subject, *domain.Subject](repo, logger, txTimeout),
		repo:        repo,
		db:          db,
	}
}

// [RU] GetByType возвращает предметы указанного типа <--->
// [ENG] GetByType returns subjects of the specified type
func (m *SubjectManager) GetByType(ctx context.Context, subjectType string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "subject_type", Operator: "=", Value: subjectType},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.Logger.Error("GetByType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "type", Value: subjectType},
		)
		return nil, fmt.Errorf("failed to get subjects by type: %w", err)
	}
	return subjects, nil
}

// [RU] SearchByName ищет предметы по названию <--->
// [ENG] SearchByName searches for subjects by name
func (m *SubjectManager) SearchByName(ctx context.Context, name string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "subject_name", Operator: "ILIKE", Value: "%" + name + "%"},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.Logger.Error("SearchByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return nil, fmt.Errorf("failed to search subjects by name: %w", err)
	}
	return subjects, nil
}

// [RU] GetByDescription ищет предметы по описанию <--->
// [ENG] GetByDescription searches for subjects by description
func (m *SubjectManager) GetByDescription(ctx context.Context, keyword string) ([]*domain.Subject, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "short_desc", Operator: "ILIKE", Value: "%" + keyword + "%"},
		},
		OrderBy: "subject_name",
	})
	if err != nil {
		m.Logger.Error("GetByDescription failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "keyword", Value: keyword},
		)
		return nil, fmt.Errorf("failed to get subjects by description: %w", err)
	}
	return subjects, nil
}

// [RU] CheckNameUnique проверяет уникальность названия предмета <--->
// [ENG] CheckNameUnique checks the uniqueness of the subject name
func (m *SubjectManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	subjects, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "subject_name", Operator: "=", Value: name},
			{Field: "subject_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(subjects) == 0, nil
}

// [RU] GetSubjectsWithPrograms возвращает предметы с привязанными программами <--->
// [ENG] GetSubjectsWithPrograms returns subjects with associated programs
func (m *SubjectManager) GetSubjectsWithPrograms(ctx context.Context, programID int) ([]*domain.Subject, error) {
	subjects, err := m.repo.GetSubjectsWithPrograms(ctx, programID)
	if err != nil {
		m.Logger.Error("GetSubjectsWithPrograms failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "program_id", Value: programID},
		)
		return nil, fmt.Errorf("failed to get subjects with programs: %w", err)
	}
	return subjects, nil
}

// [RU] GetPopularSubjects возвращает самые популярные предметы (по количеству студентов) <--->
// [ENG] GetPopularSubjects returns the most popular subjects (by student count)
func (m *SubjectManager) GetPopularSubjects(ctx context.Context, limit int) ([]*domain.Subject, error) {
	subjects, err := m.repo.GetPopularSubjects(ctx, limit)
	if err != nil {
		m.Logger.Error("GetPopularSubjects failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "limit", Value: limit},
		)
		return nil, fmt.Errorf("failed to get popular subjects: %w", err)
	}
	return subjects, nil
}

// [RU] Create создает новый предмет <--->
// [ENG] Create creates a new subject
func (m *SubjectManager) Create(ctx context.Context, subject *domain.Subject) error {
	if err := subject.Validate(); err != nil {
		m.Logger.Error("Validation failed",
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

	if err := m.Repo.Create(ctx, subject); err != nil {
		m.Logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject", Value: subject},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// [RU] Update обновляет предмет <--->
// [ENG] Update updates the subject
func (m *SubjectManager) Update(ctx context.Context, subject *domain.Subject) error {
	if err := subject.Validate(); err != nil {
		m.Logger.Error("Validation failed",
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

	if err := m.Repo.Update(ctx, subject); err != nil {
		m.Logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject", Value: subject},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// [RU] BulkCreate массово создает предметы в транзакции <--->
// [ENG] BulkCreate creates multiple subjects in a transaction
func (m *SubjectManager) BulkCreate(ctx context.Context, subjects []*domain.Subject) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.Repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, subject := range subjects {
		if err := subject.Validate(); err != nil {
			return fmt.Errorf("validation failed for subject %v: %w", subject, err)
		}

		isUnique, err := m.CheckNameUnique(ctx, subject.SubjectName, 0)
		if err != nil {
			return fmt.Errorf("uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("subject name %s already exists", subject.SubjectName)
		}

		if err := txRepo.Create(ctx, subject); err != nil {
			return fmt.Errorf("create failed for subject %v: %w", subject, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

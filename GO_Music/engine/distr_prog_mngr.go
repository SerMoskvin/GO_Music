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

type ProgrammDistributionManager struct {
	*BaseManager[int, *domain.ProgrammDistribution]
	db *sql.DB
}

func NewProgrammDistributionManager(
	repo db.Repository[*domain.ProgrammDistribution, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ProgrammDistributionManager {
	return &ProgrammDistributionManager{
		BaseManager: NewBaseManager[int, *domain.ProgrammDistribution](repo, logger, txTimeout),
		db:          db,
	}
}

func (m *ProgrammDistributionManager) GetByProgrammAndSubject(
	ctx context.Context,
	programmID int,
	subjectID int,
) (*domain.ProgrammDistribution, error) {
	distributions, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programmID},
			{Field: "subject_id", Operator: "=", Value: subjectID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByProgrammAndSubject failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm_id", Value: programmID},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return nil, fmt.Errorf("failed to get distribution: %w", err)
	}

	if len(distributions) == 0 {
		return nil, nil
	}
	return *distributions[0], nil
}

func (m *ProgrammDistributionManager) CheckExists(
	ctx context.Context,
	programmID int,
	subjectID int,
) (bool, error) {
	distributions, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programmID},
			{Field: "subject_id", Operator: "=", Value: subjectID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckExists failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm_id", Value: programmID},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return len(distributions) > 0, nil
}

func (m *ProgrammDistributionManager) GetByProgramm(
	ctx context.Context,
	programmID int,
) ([]*domain.ProgrammDistribution, error) {
	distributions, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programmID},
		},
	})
	if err != nil {
		m.logger.Error("GetByProgramm failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm_id", Value: programmID},
		)
		return nil, fmt.Errorf("failed to get distributions by programm: %w", err)
	}
	return DereferenceSlice(distributions), nil
}

func (m *ProgrammDistributionManager) GetBySubject(
	ctx context.Context,
	subjectID int,
) ([]*domain.ProgrammDistribution, error) {
	distributions, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "subject_id", Operator: "=", Value: subjectID},
		},
	})
	if err != nil {
		m.logger.Error("GetBySubject failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return nil, fmt.Errorf("failed to get distributions by subject: %w", err)
	}
	return DereferenceSlice(distributions), nil
}

func (m *ProgrammDistributionManager) Create(
	ctx context.Context,
	distribution *domain.ProgrammDistribution,
) error {
	if err := distribution.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "distribution", Value: distribution},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	exists, err := m.CheckExists(ctx, distribution.MusprogrammID, distribution.SubjectID)
	if err != nil {
		return fmt.Errorf("existence check failed: %w", err)
	}
	if exists {
		return fmt.Errorf("distribution already exists for programm %d and subject %d",
			distribution.MusprogrammID, distribution.SubjectID)
	}

	if err := m.repo.Create(ctx, &distribution); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "distribution", Value: distribution},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

func (m *ProgrammDistributionManager) BulkCreate(
	ctx context.Context,
	distributions []*domain.ProgrammDistribution,
) error {
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

	for _, distr := range distributions {
		if err := distr.Validate(); err != nil {
			return fmt.Errorf("validation failed for distribution %v: %w", distr, err)
		}

		ptrToDistr := &distr
		if err := txRepo.Create(ctx, ptrToDistr); err != nil {
			return fmt.Errorf("create failed for distribution %v: %w", distr, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

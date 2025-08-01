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

type SubjectDistributionManager struct {
	*BaseManager[int, *domain.SubjectDistribution]
	db *sql.DB
}

func NewSubjectDistributionManager(
	repo db.Repository[*domain.SubjectDistribution, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *SubjectDistributionManager {
	return &SubjectDistributionManager{
		BaseManager: NewBaseManager[int, *domain.SubjectDistribution](repo, logger, txTimeout),
		db:          db,
	}
}

func (m *SubjectDistributionManager) GetByEmployeeAndSubject(
	ctx context.Context,
	employeeID int,
	subjectID int,
) (*domain.SubjectDistribution, error) {
	distributions, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "employee_id", Operator: "=", Value: employeeID},
			{Field: "subject_id", Operator: "=", Value: subjectID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByEmployeeAndSubject failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return nil, fmt.Errorf("failed to get distribution: %w", err)
	}

	if len(distributions) == 0 {
		return nil, nil
	}
	return *distributions[0], nil
}

func (m *SubjectDistributionManager) GetByEmployee(
	ctx context.Context,
	employeeID int,
) ([]*domain.SubjectDistribution, error) {
	distributions, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "employee_id", Operator: "=", Value: employeeID},
		},
	})
	if err != nil {
		m.logger.Error("GetByEmployee failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
		)
		return nil, fmt.Errorf("failed to get distributions: %w", err)
	}
	return DereferenceSlice(distributions), nil
}

func (m *SubjectDistributionManager) GetBySubject(
	ctx context.Context,
	subjectID int,
) ([]*domain.SubjectDistribution, error) {
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
		return nil, fmt.Errorf("failed to get distributions: %w", err)
	}
	return DereferenceSlice(distributions), nil
}

func (m *SubjectDistributionManager) Create(
	ctx context.Context,
	distribution *domain.SubjectDistribution,
) error {
	if err := distribution.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "distribution", Value: distribution},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	exists, err := m.GetByEmployeeAndSubject(ctx, distribution.EmployeeID, distribution.SubjectID)
	if err != nil {
		return fmt.Errorf("existence check failed: %w", err)
	}
	if exists != nil {
		return fmt.Errorf("distribution already exists for employee %d and subject %d",
			distribution.EmployeeID, distribution.SubjectID)
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

func (m *SubjectDistributionManager) BulkCreate(
	ctx context.Context,
	distributions []*domain.SubjectDistribution,
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

func (m *SubjectDistributionManager) CheckExists(
	ctx context.Context,
	employeeID int,
	subjectID int,
) (bool, error) {
	distr, err := m.GetByEmployeeAndSubject(ctx, employeeID, subjectID)
	if err != nil {
		m.logger.Error("CheckExists failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return distr != nil, nil
}

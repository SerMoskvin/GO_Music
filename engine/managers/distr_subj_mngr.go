package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

type SubjectDistributionManager struct {
	*engine.BaseManager[int, domain.SubjectDistribution, *domain.SubjectDistribution]
	db *sql.DB
}

func NewSubjectDistributionManager(
	repo db.Repository[domain.SubjectDistribution, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *SubjectDistributionManager {
	return &SubjectDistributionManager{
		BaseManager: engine.NewBaseManager[int, domain.SubjectDistribution, *domain.SubjectDistribution](repo, logger, txTimeout),
		db:          db,
	}
}

// [RU] GetByEmployeeAndSubject возвращает распределение предмета по ID сотрудника и предмета <--->
// [ENG] GetByEmployeeAndSubject returns subject distribution by employee ID and subject ID
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
		m.Logger.Error("GetByEmployeeAndSubject failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return nil, fmt.Errorf("failed to get distribution: %w", err)
	}

	if len(distributions) == 0 {
		return nil, nil
	}
	return distributions[0], nil
}

// [RU] GetByEmployee возвращает распределения по ID сотрудника <--->
// [ENG] GetByEmployee returns distributions by employee ID
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
		m.Logger.Error("GetByEmployee failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
		)
		return nil, fmt.Errorf("failed to get distributions: %w", err)
	}
	return distributions, nil
}

// [RU] GetBySubject возвращает распределения по ID предмета <--->
// [ENG] GetBySubject returns distributions by subject ID
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
		m.Logger.Error("GetBySubject failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return nil, fmt.Errorf("failed to get distributions: %w", err)
	}
	return distributions, nil
}

// [RU] Create создает новое распределение предмета <--->
// [ENG] Create creates a new subject distribution
func (m *SubjectDistributionManager) Create(
	ctx context.Context,
	distribution *domain.SubjectDistribution,
) error {
	if err := distribution.Validate(); err != nil {
		m.Logger.Error("Validation failed",
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

	if err := m.Repo.Create(ctx, distribution); err != nil {
		m.Logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "distribution", Value: distribution},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// [RU] BulkCreate создает несколько распределений предметов в транзакции <--->
// [ENG] BulkCreate creates multiple subject distributions in a transaction
func (m *SubjectDistributionManager) BulkCreate(
	ctx context.Context,
	distributions []*domain.SubjectDistribution,
) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.Repo.WithTx(tx)

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

		ptrToDistr := distr
		if err := txRepo.Create(ctx, ptrToDistr); err != nil {
			return fmt.Errorf("create failed for distribution %v: %w", distr, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// [RU] CheckExists проверяет существование распределения предмета <--->
// [ENG] CheckExists checks if subject distribution exists
func (m *SubjectDistributionManager) CheckExists(
	ctx context.Context,
	employeeID int,
	subjectID int,
) (bool, error) {
	distr, err := m.GetByEmployeeAndSubject(ctx, employeeID, subjectID)
	if err != nil {
		m.Logger.Error("CheckExists failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee_id", Value: employeeID},
			logger.Field{Key: "subject_id", Value: subjectID},
		)
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return distr != nil, nil
}

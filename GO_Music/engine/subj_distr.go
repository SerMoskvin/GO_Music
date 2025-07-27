package engine

import (
	"context"
	"fmt"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

type SubjectDistributionManager struct {
	*BaseManager[domain.SubjectDistribution, *domain.SubjectDistribution]
}

func NewSubjectDistributionManager(
	repo Repository[domain.SubjectDistribution, *domain.SubjectDistribution],
	logger *logger.LevelLogger,
) *SubjectDistributionManager {
	return &SubjectDistributionManager{
		BaseManager: &BaseManager[domain.SubjectDistribution, *domain.SubjectDistribution]{
			repo:   repo,
			logger: logger,
		},
	}
}

func (m *SubjectDistributionManager) GetByEmployeeAndSubject(
	ctx context.Context,
	employeeID int,
	subjectID int,
) (*domain.SubjectDistribution, error) {
	distributions, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return distributions[0], nil
}

func (m *SubjectDistributionManager) GetByEmployee(
	ctx context.Context,
	employeeID int,
) ([]*domain.SubjectDistribution, error) {
	distributions, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return distributions, nil
}

func (m *SubjectDistributionManager) GetBySubject(
	ctx context.Context,
	subjectID int,
) ([]*domain.SubjectDistribution, error) {
	distributions, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return distributions, nil
}

package engine

import (
	"context"
	"fmt"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

type ProgrammDistributionManager struct {
	*BaseManager[domain.ProgrammDistribution, *domain.ProgrammDistribution]
}

func NewProgrammDistributionManager(
	repo Repository[domain.ProgrammDistribution, *domain.ProgrammDistribution],
	logger *logger.LevelLogger,
) *ProgrammDistributionManager {
	return &ProgrammDistributionManager{
		BaseManager: &BaseManager[domain.ProgrammDistribution, *domain.ProgrammDistribution]{
			repo:   repo,
			logger: logger,
		},
	}
}

func (m *ProgrammDistributionManager) GetByProgrammAndSubject(
	ctx context.Context,
	programmID int,
	subjectID int,
) (*domain.ProgrammDistribution, error) {
	distributions, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return distributions[0], nil
}

func (m *ProgrammDistributionManager) CheckExists(
	ctx context.Context,
	programmID int,
	subjectID int,
) (bool, error) {
	distributions, err := m.List(ctx, Filter{
		Conditions: []Condition{
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

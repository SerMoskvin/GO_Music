package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// StudyGroupManager реализует бизнес-логику для учебных групп
type StudyGroupManager struct {
	*BaseManager[domain.StudyGroup, *domain.StudyGroup]
}

func NewStudyGroupManager(
	repo Repository[domain.StudyGroup, *domain.StudyGroup],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudyGroupManager {
	return &StudyGroupManager{
		BaseManager: NewBaseManager[domain.StudyGroup](repo, logger, txTimeout),
	}
}

// GetByProgram возвращает группы по программе обучения
func (m *StudyGroupManager) GetByProgram(ctx context.Context, programID int) ([]*domain.StudyGroup, error) {
	groups, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programID},
		},
		OrderBy: "study_year DESC, group_name",
	})
	if err != nil {
		m.logger.Error("GetByProgram failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "program_id", Value: programID},
		)
		return nil, fmt.Errorf("failed to get groups by program: %w", err)
	}
	return groups, nil
}

// GetByName возвращает группу по названию (точное совпадение)
func (m *StudyGroupManager) GetByName(ctx context.Context, name string) (*domain.StudyGroup, error) {
	groups, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "group_name", Operator: "=", Value: name},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return nil, fmt.Errorf("failed to get group by name: %w", err)
	}

	if len(groups) == 0 {
		return nil, nil
	}
	return groups[0], nil
}

// GetByYear возвращает группы по учебному году
func (m *StudyGroupManager) GetByYear(ctx context.Context, year int) ([]*domain.StudyGroup, error) {
	groups, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "study_year", Operator: "=", Value: year},
		},
		OrderBy: "group_name",
	})
	if err != nil {
		m.logger.Error("GetByYear failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "year", Value: year},
		)
		return nil, fmt.Errorf("failed to get groups by year: %w", err)
	}
	return groups, nil
}

// CheckNameUnique проверяет уникальность названия группы
func (m *StudyGroupManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	groups, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "group_name", Operator: "=", Value: name},
			{Field: "group_id", Operator: "!=", Value: excludeID},
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
	return len(groups) == 0, nil
}

// UpdateStudentCount обновляет количество студентов в группе
func (m *StudyGroupManager) UpdateStudentCount(ctx context.Context, groupID int, newCount int) error {
	group, err := m.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	group.NumberOfStudents = newCount
	return m.Update(ctx, group)
}

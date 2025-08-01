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

// StudyGroupManager реализует бизнес-логику для учебных групп
type StudyGroupManager struct {
	*BaseManager[int, *domain.StudyGroup]
	db *sql.DB
}

func NewStudyGroupManager(
	repo db.Repository[*domain.StudyGroup, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudyGroupManager {
	return &StudyGroupManager{
		BaseManager: NewBaseManager[int, *domain.StudyGroup](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByProgram возвращает группы по программе обучения
func (m *StudyGroupManager) GetByProgram(ctx context.Context, programID int) ([]*domain.StudyGroup, error) {
	groups, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(groups), nil
}

// GetByName возвращает группу по названию (точное совпадение)
func (m *StudyGroupManager) GetByName(ctx context.Context, name string) (*domain.StudyGroup, error) {
	groups, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return *groups[0], nil
}

// GetByYear возвращает группы по учебному году
func (m *StudyGroupManager) GetByYear(ctx context.Context, year int) ([]*domain.StudyGroup, error) {
	groups, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(groups), nil
}

// CheckNameUnique проверяет уникальность названия группы
func (m *StudyGroupManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	groups, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "group_name", Operator: "=", Value: name},
			{Field: "group_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "exclude_id", Value: excludeID},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(groups) == 0, nil
}

// UpdateStudentCount обновляет количество студентов в группе
func (m *StudyGroupManager) UpdateStudentCount(ctx context.Context, groupID int, newCount int) error {
	groupPtr, err := m.GetByID(ctx, groupID)
	if err != nil {
		m.logger.Error("UpdateStudentCount failed to get group",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group_id", Value: groupID},
		)
		return fmt.Errorf("failed to get group: %w", err)
	}
	if groupPtr == nil {
		m.logger.Error("Group not found",
			logger.Field{Key: "group_id", Value: groupID},
		)
		return fmt.Errorf("group not found")
	}

	// Разыменовываем указатель для изменения значения
	group := *groupPtr
	group.NumberOfStudents = newCount

	if err := m.repo.Update(ctx, &group); err != nil {
		m.logger.Error("UpdateStudentCount failed to update",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group", Value: group},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// Create создает новую учебную группу
func (m *StudyGroupManager) Create(ctx context.Context, group *domain.StudyGroup) error {
	if err := group.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group", Value: group},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, group.GroupName, 0)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("group name %s already exists", group.GroupName)
	}

	if err := m.repo.Create(ctx, &group); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group", Value: group},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// Update обновляет данные учебной группы
func (m *StudyGroupManager) Update(ctx context.Context, group *domain.StudyGroup) error {
	if err := group.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group", Value: group},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, group.GroupName, group.GroupID)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("group name %s already exists", group.GroupName)
	}

	if err := m.repo.Update(ctx, &group); err != nil {
		m.logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group", Value: group},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// BulkCreate массово создает учебные группы в транзакции
func (m *StudyGroupManager) BulkCreate(ctx context.Context, groups []*domain.StudyGroup) error {
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

	for _, gr := range groups {
		if err := gr.Validate(); err != nil {
			return fmt.Errorf("validation failed for group %v: %w", gr, err)
		}

		isUnique, err := m.CheckNameUnique(ctx, gr.GroupName, 0)
		if err != nil {
			return fmt.Errorf("name uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("group name %s already exists", gr.GroupName)
		}

		ptrToGroup := &gr
		if err := txRepo.Create(ctx, ptrToGroup); err != nil {
			return fmt.Errorf("create failed for group %v: %w", gr, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

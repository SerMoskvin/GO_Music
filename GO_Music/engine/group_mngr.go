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

type StudyGroupManager struct {
	*BaseManager[int, domain.StudyGroup, *domain.StudyGroup]
	db *sql.DB
}

func NewStudyGroupManager(
	repo db.Repository[domain.StudyGroup, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudyGroupManager {
	return &StudyGroupManager{
		BaseManager: NewBaseManager[int, domain.StudyGroup, *domain.StudyGroup](repo, logger, txTimeout),
		db:          db,
	}
}

// [RU] GetByProgram возвращает группы по программе обучения <--->
// [ENG] GetByProgram returns groups by training program
func (m *StudyGroupManager) GetByProgram(ctx context.Context, programID int) ([]*domain.StudyGroup, error) {
	groups, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "mus_programm_id", Operator: "=", Value: programID},
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

// [RU] GetByName возвращает группу по названию (точное совпадение) <--->
// [ENG] GetByName returns a group by name (exact match)
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
	return groups[0], nil
}

// [RU] GetByYear возвращает группы по учебному году <--->
// [ENG] GetByYear returns groups by academic year
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
	return groups, nil
}

// [RU] CheckNameUnique проверяет уникальность названия группы <--->
// [ENG] CheckNameUnique checks the uniqueness of the group name
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

// [RU] UpdateStudentCount обновляет количество студентов в группе <--->
// [ENG] UpdateStudentCount updates the number of students in the group
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

	group := *groupPtr
	group.NumberOfStudents = newCount

	if err := m.BaseManager.Update(ctx, &group); err != nil {
		m.logger.Error("UpdateStudentCount failed to update",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group", Value: group},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// [RU] Create создает новую учебную группу <--->
// [ENG] Create creates a new study group
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

	return m.BaseManager.Create(ctx, group)
}

// [RU] Update обновляет данные учебной группы <--->
// [ENG] Update updates the study group data
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

	return m.BaseManager.Update(ctx, group)
}

// [RU] BulkCreate массово создает учебные группы в транзакции <--->
// [ENG] BulkCreate creates multiple study groups in a transaction
func (m *StudyGroupManager) BulkCreate(ctx context.Context, groups []*domain.StudyGroup) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
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

		if err := txRepo.Create(ctx, gr); err != nil {
			return fmt.Errorf("create failed for group %v: %w", gr, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

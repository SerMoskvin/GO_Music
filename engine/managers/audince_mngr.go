package managers

import (
	"context"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

// AudienceManager реализует бизнес-логику для аудиторий
type AudienceManager struct {
	*engine.BaseManager[int, domain.Audience, *domain.Audience]
}

// NewAudienceManager создает новый экземпляр AudienceManager
func NewAudienceManager(
	repo db.Repository[domain.Audience, int],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *AudienceManager {
	return &AudienceManager{
		BaseManager: engine.NewBaseManager[int, domain.Audience, *domain.Audience](repo, logger, txTimeout),
	}
}

// [RU] GetByNumber возвращает аудиторию по номеру <--->
// [ENG] GetByNumber returns audience by number
func (m *AudienceManager) GetByNumber(ctx context.Context, number string) (*domain.Audience, error) {
	audiences, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "audin_number", Operator: "=", Value: number},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("GetByNumber failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "number", Value: number},
		)
		return nil, fmt.Errorf("get by number failed: %w", err)
	}

	if len(audiences) == 0 {
		return nil, nil
	}
	return audiences[0], nil
}

// [RU] ListByCapacity возвращает аудитории с вместимостью >= minCapacity <--->
// [ENG] ListByCapacity returns audiences with capacity >= minCapacity
func (m *AudienceManager) ListByCapacity(ctx context.Context, minCapacity int) ([]*domain.Audience, error) {
	audiences, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "capacity", Operator: ">=", Value: minCapacity},
		},
		OrderBy: "capacity DESC",
	})
	if err != nil {
		m.Logger.Error("ListByCapacity failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "minCapacity", Value: minCapacity},
		)
		return nil, fmt.Errorf("list by capacity failed: %w", err)
	}
	return audiences, nil
}

// [RU] CheckNumberUnique проверяет уникальность номера аудитории <--->
// [ENG] CheckNumberUnique checks audience number uniqueness
func (m *AudienceManager) CheckNumberUnique(ctx context.Context, number string, excludeID int) (bool, error) {
	audiences, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "audin_number", Operator: "=", Value: number},
			{Field: "audience_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckNumberUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "number", Value: number},
			logger.Field{Key: "excludeID", Value: excludeID},
		)
		return false, fmt.Errorf("check number unique failed: %w", err)
	}
	return len(audiences) == 0, nil
}

// [RU] Create создает новую аудиторию с проверкой уникальности номера <--->
// [ENG] Create creates a new audience with uniqueness check for number
func (m *AudienceManager) Create(ctx context.Context, audience *domain.Audience) error {
	if err := audience.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience", Value: audience},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNumberUnique(ctx, audience.AudinNumber, 0)
	if err != nil {
		return fmt.Errorf("uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("audience number %s already exists", audience.AudinNumber)
	}

	if err := m.Repo.Create(ctx, audience); err != nil {
		m.Logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience", Value: audience},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// [RU] Update обновляет аудиторию с проверкой уникальности номера <--->
// [ENG] Update updates audience with uniqueness check for number
func (m *AudienceManager) Update(ctx context.Context, audience *domain.Audience) error {
	if err := audience.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience", Value: audience},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNumberUnique(ctx, audience.AudinNumber, audience.AudienceID)
	if err != nil {
		return fmt.Errorf("uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("audience number %s already exists", audience.AudinNumber)
	}

	if err := m.Repo.Update(ctx, audience); err != nil {
		m.Logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience", Value: audience},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

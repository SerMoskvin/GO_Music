package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// AudienceManager расширяет BaseManager для аудиторий
type AudienceManager struct {
	*BaseManager[domain.Audience, *domain.Audience]
}

func NewAudienceManager(
	repo Repository[domain.Audience, *domain.Audience],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *AudienceManager {
	return &AudienceManager{
		BaseManager: NewBaseManager[domain.Audience](repo, logger, txTimeout),
	}
}

// GetByNumber возвращает аудиторию по номеру
func (m *AudienceManager) GetByNumber(ctx context.Context, number string) (*domain.Audience, error) {
	audiences, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "audin_number", Operator: "=", Value: number},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByNumber failed",
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

// ListByCapacity возвращает аудитории с вместимостью >= minCapacity
func (m *AudienceManager) ListByCapacity(ctx context.Context, minCapacity int) ([]*domain.Audience, error) {
	return m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "capacity", Operator: ">=", Value: minCapacity},
		},
		OrderBy: "capacity DESC",
	})
}

// CheckNumberUnique проверяет уникальность номера аудитории
func (m *AudienceManager) CheckNumberUnique(ctx context.Context, number string, excludeID int) (bool, error) {
	audiences, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "audin_number", Operator: "=", Value: number},
			{Field: "audience_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		return false, err
	}
	return len(audiences) == 0, nil
}

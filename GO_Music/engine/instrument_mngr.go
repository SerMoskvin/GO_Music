package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// InstrumentManager реализует бизнес-логику для музыкальных инструментов
type InstrumentManager struct {
	*BaseManager[domain.Instrument, *domain.Instrument]
}

func NewInstrumentManager(
	repo Repository[domain.Instrument, *domain.Instrument],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *InstrumentManager {
	return &InstrumentManager{
		BaseManager: NewBaseManager[domain.Instrument](repo, logger, txTimeout),
	}
}

// GetByAudience возвращает инструменты в указанной аудитории
func (m *InstrumentManager) GetByAudience(ctx context.Context, audienceID int) ([]*domain.Instrument, error) {
	instruments, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "audience_id", Operator: "=", Value: audienceID},
		},
		OrderBy: "name",
	})
	if err != nil {
		m.logger.Error("GetByAudience failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience_id", Value: audienceID},
		)
		return nil, fmt.Errorf("failed to get instruments by audience: %w", err)
	}
	return instruments, nil
}

// GetByType возвращает инструменты указанного типа
func (m *InstrumentManager) GetByType(ctx context.Context, instrType string) ([]*domain.Instrument, error) {
	instruments, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "instr_type", Operator: "=", Value: instrType},
		},
		OrderBy: "name",
	})
	if err != nil {
		m.logger.Error("GetByType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "type", Value: instrType},
		)
		return nil, fmt.Errorf("failed to get instruments by type: %w", err)
	}
	return instruments, nil
}

// GetByName возвращает инструмент по точному названию
func (m *InstrumentManager) GetByName(ctx context.Context, name string) (*domain.Instrument, error) {
	instruments, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "name", Operator: "=", Value: name},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return nil, fmt.Errorf("failed to get instrument by name: %w", err)
	}

	if len(instruments) == 0 {
		return nil, nil
	}
	return instruments[0], nil
}

// CheckNameUnique проверяет уникальность названия инструмента
func (m *InstrumentManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	instruments, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "name", Operator: "=", Value: name},
			{Field: "instrument_id", Operator: "!=", Value: excludeID},
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
	return len(instruments) == 0, nil
}

// UpdateCondition обновляет состояние инструмента
func (m *InstrumentManager) UpdateCondition(ctx context.Context, instrumentID int, newCondition string) error {
	instrument, err := m.GetByID(ctx, instrumentID)
	if err != nil {
		return fmt.Errorf("failed to get instrument: %w", err)
	}
	if instrument == nil {
		return fmt.Errorf("instrument not found")
	}

	instrument.Condition = newCondition
	return m.Update(ctx, instrument)
}

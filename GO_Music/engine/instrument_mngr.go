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

// InstrumentManager реализует бизнес-логику для музыкальных инструментов
type InstrumentManager struct {
	*BaseManager[int, *domain.Instrument]
	db *sql.DB
}

func NewInstrumentManager(
	repo db.Repository[*domain.Instrument, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *InstrumentManager {
	return &InstrumentManager{
		BaseManager: NewBaseManager[int, *domain.Instrument](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByAudience возвращает инструменты в указанной аудитории
func (m *InstrumentManager) GetByAudience(ctx context.Context, audienceID int) ([]*domain.Instrument, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(instruments), nil
}

// GetByType возвращает инструменты указанного типа
func (m *InstrumentManager) GetByType(ctx context.Context, instrType string) ([]*domain.Instrument, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return DereferenceSlice(instruments), nil
}

// GetByName возвращает инструмент по точному названию
func (m *InstrumentManager) GetByName(ctx context.Context, name string) (*domain.Instrument, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
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
	return *instruments[0], nil
}

// CheckNameUnique проверяет уникальность названия инструмента
func (m *InstrumentManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "name", Operator: "=", Value: name},
			{Field: "instrument_id", Operator: "!=", Value: excludeID},
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
	return len(instruments) == 0, nil
}

// UpdateCondition обновляет состояние инструмента
func (m *InstrumentManager) UpdateCondition(ctx context.Context, instrumentID int, newCondition string) error {
	instrumentPtr, err := m.GetByID(ctx, instrumentID)
	if err != nil {
		m.logger.Error("UpdateCondition failed to get instrument",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument_id", Value: instrumentID},
		)
		return fmt.Errorf("failed to get instrument: %w", err)
	}
	if instrumentPtr == nil {
		m.logger.Error("Instrument not found",
			logger.Field{Key: "instrument_id", Value: instrumentID},
		)
		return fmt.Errorf("instrument not found")
	}

	instrument := *instrumentPtr
	instrument.Condition = newCondition

	if err := m.repo.Update(ctx, &instrument); err != nil {
		m.logger.Error("UpdateCondition failed to update",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// Create создает новый инструмент
func (m *InstrumentManager) Create(ctx context.Context, instrument *domain.Instrument) error {
	if err := instrument.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Проверяем уникальность названия инструмента
	isUnique, err := m.CheckNameUnique(ctx, instrument.Name, 0)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("instrument name %s already exists", instrument.Name)
	}

	if err := m.repo.Create(ctx, &instrument); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// Update обновляет данные инструмента
func (m *InstrumentManager) Update(ctx context.Context, instrument *domain.Instrument) error {
	if err := instrument.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Проверяем уникальность названия инструмента
	isUnique, err := m.CheckNameUnique(ctx, instrument.Name, instrument.InstrumentID)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("instrument name %s already exists", instrument.Name)
	}

	if err := m.repo.Update(ctx, &instrument); err != nil {
		m.logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// BulkCreate массово создает инструменты в транзакции
func (m *InstrumentManager) BulkCreate(ctx context.Context, instruments []*domain.Instrument) error {
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

	for _, instr := range instruments {
		if err := instr.Validate(); err != nil {
			return fmt.Errorf("validation failed for instrument %v: %w", instr, err)
		}

		// Проверяем уникальность названия инструмента
		isUnique, err := m.CheckNameUnique(ctx, instr.Name, 0)
		if err != nil {
			return fmt.Errorf("name uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("instrument name %s already exists", instr.Name)
		}

		ptrToInstr := &instr
		if err := txRepo.Create(ctx, ptrToInstr); err != nil {
			return fmt.Errorf("create failed for instrument %v: %w", instr, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

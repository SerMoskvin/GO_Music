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

type InstrumentManager struct {
	*engine.BaseManager[int, domain.Instrument, *domain.Instrument]
	db *sql.DB
}

func NewInstrumentManager(
	repo db.Repository[domain.Instrument, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *InstrumentManager {
	return &InstrumentManager{
		BaseManager: engine.NewBaseManager[int, domain.Instrument, *domain.Instrument](repo, logger, txTimeout),
		db:          db,
	}
}

// [RU] GetByAudience возвращает инструменты в указанной аудитории <--->
// [ENG] GetByAudience returns instruments in the specified audience
func (m *InstrumentManager) GetByAudience(ctx context.Context, audienceID int) ([]*domain.Instrument, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "audience_id", Operator: "=", Value: audienceID},
		},
		OrderBy: "name",
	})
	if err != nil {
		m.Logger.Error("GetByAudience failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "audience_id", Value: audienceID},
		)
		return nil, fmt.Errorf("failed to get instruments by audience: %w", err)
	}
	return instruments, nil
}

// [RU] GetByType возвращает инструменты указанного типа <--->
// [ENG] GetByType returns instruments of the specified type
func (m *InstrumentManager) GetByType(ctx context.Context, instrType string) ([]*domain.Instrument, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "instr_type", Operator: "=", Value: instrType},
		},
		OrderBy: "name",
	})
	if err != nil {
		m.Logger.Error("GetByType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "type", Value: instrType},
		)
		return nil, fmt.Errorf("failed to get instruments by type: %w", err)
	}
	return instruments, nil
}

// [RU] GetByName возвращает инструмент по точному названию <--->
// [ENG] GetByName returns an instrument by exact name
func (m *InstrumentManager) GetByName(ctx context.Context, name string) (*domain.Instrument, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "name", Operator: "=", Value: name},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("GetByName failed",
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

// [RU] CheckNameUnique проверяет уникальность названия инструмента (исключая указанный ID) <--->
// [ENG] CheckNameUnique checks the uniqueness of the instrument name excluding the given ID
func (m *InstrumentManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	instruments, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "name", Operator: "=", Value: name},
			{Field: "instrument_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "exclude_id", Value: excludeID},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(instruments) == 0, nil
}

// [RU] Create создает новый инструмент с проверкой уникальности названия <--->
// [ENG] Create creates a new instrument with name uniqueness check
func (m *InstrumentManager) Create(ctx context.Context, instrument *domain.Instrument) error {
	if err := instrument.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, instrument.Name, 0)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("instrument name %s already exists", instrument.Name)
	}

	return m.BaseManager.Create(ctx, instrument)
}

// [RU] Update обновляет данные инструмента с проверкой уникальности названия <--->
// [ENG] Update updates instrument data with name uniqueness check
func (m *InstrumentManager) Update(ctx context.Context, instrument *domain.Instrument) error {
	if err := instrument.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, instrument.Name, instrument.InstrumentID)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("instrument name %s already exists", instrument.Name)
	}

	return m.BaseManager.Update(ctx, instrument)
}

// [RU] UpdateCondition обновляет состояние инструмента <--->
// [ENG] UpdateCondition updates the condition of the instrument
func (m *InstrumentManager) UpdateCondition(ctx context.Context, instrumentID int, newCondition string) error {
	instrument, err := m.GetByID(ctx, instrumentID)
	if err != nil {
		m.Logger.Error("UpdateCondition failed to get instrument",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument_id", Value: instrumentID},
		)
		return fmt.Errorf("failed to get instrument: %w", err)
	}
	if instrument == nil {
		m.Logger.Error("Instrument not found",
			logger.Field{Key: "instrument_id", Value: instrumentID},
		)
		return fmt.Errorf("instrument not found")
	}

	instrument.Condition = newCondition

	return m.BaseManager.Update(ctx, instrument)
}

// [RU] BulkCreate массово создает инструменты в транзакции <--->
// [ENG] BulkCreate creates multiple instruments in a transaction
func (m *InstrumentManager) BulkCreate(ctx context.Context, instruments []*domain.Instrument) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.Repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, instr := range instruments {
		if err := instr.Validate(); err != nil {
			return fmt.Errorf("validation failed for instrument %v: %w", instr, err)
		}

		isUnique, err := m.CheckNameUnique(ctx, instr.Name, 0)
		if err != nil {
			return fmt.Errorf("name uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("instrument name %s already exists", instr.Name)
		}

		if err := txRepo.Create(ctx, instr); err != nil {
			return fmt.Errorf("create failed for instrument %v: %w", instr, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

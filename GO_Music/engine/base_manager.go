package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// [RU] TxProvider интерфейс для работы с транзакциями <--->
// [ENG] Txprovider - interface for work with transaction
type TxProvider interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// [RU] BaseManager - базовая реализация CRUD операций <--->
// [ENG] BaseManager - basic realisation for CRUD operations
type BaseManager[ID comparable, T any, PT interface {
	*T
	domain.Entity[ID]
}] struct {
	Repo      db.Repository[T, ID]
	Logger    *logger.LevelLogger
	txTimeout time.Duration
}

// Конструктор менеджера
func NewBaseManager[ID comparable, T any, PT interface {
	*T
	domain.Entity[ID]
}](
	repo db.Repository[T, ID],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *BaseManager[ID, T, PT] {
	if logger == nil {
		panic("logger is required")
	}
	return &BaseManager[ID, T, PT]{
		Repo:      repo,
		Logger:    logger,
		txTimeout: txTimeout,
	}
}

func (m *BaseManager[ID, T, PT]) Create(ctx context.Context, entity PT) error {
	if err := entity.Validate(); err != nil {
		m.Logger.Error("Validation failed", logger.Error(err))
		return fmt.Errorf("validation error: %w", err)
	}

	if err := m.Repo.Create(ctx, entity); err != nil {
		m.Logger.Error("Create failed", logger.Error(err))
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

func (m *BaseManager[ID, T, PT]) Update(ctx context.Context, entity PT) error {
	if entity.GetID() == *new(ID) {
		return errors.New("ID is required")
	}

	if err := entity.Validate(); err != nil {
		m.Logger.Error("Validation failed", logger.Error(err))
		return fmt.Errorf("validation error: %w", err)
	}

	if err := m.Repo.Update(ctx, entity); err != nil {
		m.Logger.Error("Update failed", logger.Error(err))
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (m *BaseManager[ID, T, PT]) Delete(ctx context.Context, id ID) error {
	var zeroID ID
	if id == zeroID {
		return errors.New("ID is required")
	}

	if err := m.Repo.Delete(ctx, id); err != nil {
		m.Logger.Error("Delete failed", logger.Error(err), logger.Any("id", id))
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

func (m *BaseManager[ID, T, PT]) GetByID(ctx context.Context, id ID) (PT, error) {
	var zeroID ID
	if id == zeroID {
		return nil, errors.New("ID is required")
	}

	entity, err := m.Repo.GetByID(ctx, id)
	if err != nil {
		m.Logger.Error("GetByID failed", logger.Error(err), logger.Any("id", id))
		return nil, fmt.Errorf("get failed: %w", err)
	}
	return entity, nil
}

func (m *BaseManager[ID, T, PT]) GetByIDs(ctx context.Context, ids []ID) ([]*T, error) {
	if len(ids) == 0 {
		return nil, errors.New("at least one ID required")
	}

	entities, err := m.Repo.GetByIDs(ctx, ids)
	if err != nil {
		m.Logger.Error("GetByIDs failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "ids", Value: ids},
		)
		return nil, fmt.Errorf("get multiple failed: %w", err)
	}
	return entities, nil
}

func (m *BaseManager[ID, T, PT]) List(ctx context.Context, filter db.Filter) ([]*T, error) {
	entities, err := m.Repo.List(ctx, filter)
	if err != nil {
		m.Logger.Error("List failed", logger.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("list failed: %w", err)
	}
	return entities, nil
}

func (m *BaseManager[ID, T, PT]) Count(ctx context.Context, filter db.Filter) (int, error) {
	count, err := m.Repo.Count(ctx, filter)
	if err != nil {
		m.Logger.Error("Count failed", logger.Error(err))
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return count, nil
}

func (m *BaseManager[ID, T, PT]) Exists(ctx context.Context, id ID) (bool, error) {
	exists, err := m.Repo.Exists(ctx, id)
	if err != nil {
		m.Logger.Error("Exists check failed", logger.Error(err), logger.Any("id", id))
		return false, fmt.Errorf("exists check failed: %w", err)
	}
	return exists, nil
}

func (m *BaseManager[ID, T, PT]) ExecuteInTx(
	ctx context.Context,
	txProvider TxProvider,
	ops func(repo db.Repository[T, ID]) error,
) error {
	txCtx, cancel := context.WithTimeout(ctx, m.txTimeout)
	defer cancel()

	tx, err := txProvider.BeginTx(txCtx, nil)
	if err != nil {
		m.Logger.Error("BeginTx failed", logger.Error(err))
		return fmt.Errorf("begin tx failed: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := ops(m.Repo.WithTx(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		m.Logger.Error("Commit failed", logger.Error(err))
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

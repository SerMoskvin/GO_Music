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

// TxProvider интерфейс для работы с транзакциями
type TxProvider interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// BaseManager - базовая реализация CRUD операций
type BaseManager[ID comparable, T domain.Entity[ID]] struct {
	repo      db.Repository[T, ID]
	logger    *logger.LevelLogger
	txTimeout time.Duration
}

// Конструктор менеджера
func NewBaseManager[ID comparable, T domain.Entity[ID]](
	repo db.Repository[T, ID],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *BaseManager[ID, T] {
	if logger == nil {
		panic("logger is required")
	}
	return &BaseManager[ID, T]{
		repo:      repo,
		logger:    logger,
		txTimeout: txTimeout,
	}
}

func (m *BaseManager[ID, T]) Create(ctx context.Context, entity *T) error {
	if err := (*entity).Validate(); err != nil {
		m.logger.Error("Validation failed", logger.Field{Key: "error", Value: err})
		return fmt.Errorf("validation error: %w", err)
	}

	if err := m.repo.Create(ctx, entity); err != nil {
		m.logger.Error("Create failed", logger.Field{Key: "error", Value: err})
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

func (m *BaseManager[ID, T]) Update(ctx context.Context, entity *T) error {
	if (*entity).GetID() == *new(ID) {
		return errors.New("ID is required")
	}

	if err := (*entity).Validate(); err != nil {
		m.logger.Error("Validation failed", logger.Error(err))
		return fmt.Errorf("validation error: %w", err)
	}

	if err := m.repo.Update(ctx, entity); err != nil {
		m.logger.Error("Update failed", logger.Error(err))
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (m *BaseManager[ID, T]) Delete(ctx context.Context, id ID) error {
	var zeroID ID
	if id == zeroID {
		return errors.New("ID is required")
	}

	if err := m.repo.Delete(ctx, id); err != nil {
		m.logger.Error("Delete failed", logger.Error(err), logger.Any("id", id))
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

func (m *BaseManager[ID, T]) GetByID(ctx context.Context, id ID) (*T, error) {
	var zeroID ID
	if id == zeroID {
		return nil, errors.New("ID is required")
	}

	entity, err := m.repo.GetByID(ctx, id)
	if err != nil {
		m.logger.Error("GetByID failed", logger.Error(err), logger.Any("id", id))
		return nil, fmt.Errorf("get failed: %w", err)
	}
	return entity, nil
}

func (m *BaseManager[ID, T]) GetByIDs(ctx context.Context, ids []ID) ([]*T, error) {
	if len(ids) == 0 {
		return nil, errors.New("at least one ID required")
	}

	entities, err := m.repo.GetByIDs(ctx, ids)
	if err != nil {
		m.logger.Error("GetByIDs failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "ids", Value: ids},
		)
		return nil, fmt.Errorf("get multiple failed: %w", err)
	}
	return entities, nil
}

func (m *BaseManager[ID, T]) List(ctx context.Context, filter db.Filter) ([]*T, error) {
	entities, err := m.repo.List(ctx, filter)
	if err != nil {
		m.logger.Error("List failed", logger.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("list failed: %w", err)
	}
	return entities, nil
}

func (m *BaseManager[ID, T]) Count(ctx context.Context, filter db.Filter) (int, error) {
	count, err := m.repo.Count(ctx, filter)
	if err != nil {
		m.logger.Error("Count failed", logger.Error(err))
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return count, nil
}

func (m *BaseManager[ID, T]) Exists(ctx context.Context, id ID) (bool, error) {
	exists, err := m.repo.Exists(ctx, id)
	if err != nil {
		m.logger.Error("Exists check failed", logger.Error(err), logger.Any("id", id))
		return false, fmt.Errorf("exists check failed: %w", err)
	}
	return exists, nil
}

func (m *BaseManager[ID, T]) ExecuteInTx(
	ctx context.Context,
	txProvider TxProvider,
	ops func(repo db.Repository[T, ID]) error,
) error {
	txCtx, cancel := context.WithTimeout(ctx, m.txTimeout)
	defer cancel()

	tx, err := txProvider.BeginTx(txCtx, nil)
	if err != nil {
		m.logger.Error("BeginTx failed", logger.Error(err))
		return fmt.Errorf("begin tx failed: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := ops(m.repo.WithTx(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		m.logger.Error("Commit failed", logger.Error(err))
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

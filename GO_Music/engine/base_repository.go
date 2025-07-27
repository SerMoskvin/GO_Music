package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/SerMoskvin/logger"
)

// Entity определяет базовые методы для всех сущностей
type Entity interface {
	GetID() int
	SetID(id int)
	Validate() error
}

// PointerEntity - интерфейс для работы с указателями на сущности
type PointerEntity[T any] interface {
	*T
	Entity
}

// TxProvider интерфейс для работы с транзакциями
type TxProvider interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Filter для пагинации и фильтрации
type Filter struct {
	Limit      int
	Offset     int
	OrderBy    string
	Preloads   []string
	Search     string
	Conditions []Condition
}

type Condition struct {
	Field    string
	Operator string
	Value    interface{}
}

// Repository - полный интерфейс репозитория
type Repository[T any, PT PointerEntity[T]] interface {
	Create(ctx context.Context, entity PT) error
	Update(ctx context.Context, entity PT) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (PT, error)
	GetByIDs(ctx context.Context, ids []int) ([]PT, error)
	List(ctx context.Context, filter Filter) ([]PT, error)
	Count(ctx context.Context, filter Filter) (int, error)
	Exists(ctx context.Context, id int) (bool, error)
	WithTx(tx *sql.Tx) Repository[T, PT]
}

// BaseManager - базовая реализация CRUD операций
type BaseManager[T any, PT PointerEntity[T]] struct {
	repo      Repository[T, PT]
	logger    *logger.LevelLogger
	txTimeout time.Duration
}

func NewBaseManager[T any, PT PointerEntity[T]](
	repo Repository[T, PT],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *BaseManager[T, PT] {
	if logger == nil {
		panic("logger is required")
	}
	return &BaseManager[T, PT]{
		repo:      repo,
		logger:    logger,
		txTimeout: txTimeout,
	}
}

func (m *BaseManager[T, PT]) Create(ctx context.Context, entity PT) error {
	if err := entity.Validate(); err != nil {
		m.logger.Error("Validation failed", logger.Field{Key: "error", Value: err})
		return fmt.Errorf("validation error: %w", err)
	}

	if err := m.repo.Create(ctx, entity); err != nil {
		m.logger.Error("Create failed", logger.Field{Key: "error", Value: err})
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

func (m *BaseManager[T, PT]) Update(ctx context.Context, entity PT) error {
	if entity.GetID() == 0 {
		return errors.New("ID is required")
	}

	if err := entity.Validate(); err != nil {
		m.logger.Error("Validation failed", logger.Error(err))
		return fmt.Errorf("validation error: %w", err)
	}

	if err := m.repo.Update(ctx, entity); err != nil {
		m.logger.Error("Update failed", logger.Error(err))
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (m *BaseManager[T, PT]) Delete(ctx context.Context, id int) error {
	if id == 0 {
		return errors.New("ID is required")
	}

	if err := m.repo.Delete(ctx, id); err != nil {
		m.logger.Error("Delete failed", logger.Error(err), logger.Int("id", id))
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

func (m *BaseManager[T, PT]) GetByID(ctx context.Context, id int) (PT, error) {
	if id == 0 {
		return nil, errors.New("ID is required")
	}

	entity, err := m.repo.GetByID(ctx, id)
	if err != nil {
		m.logger.Error("GetByID failed", logger.Error(err), logger.Int("id", id))
		return nil, fmt.Errorf("get failed: %w", err)
	}
	return entity, nil
}

func (m *BaseManager[T, PT]) GetByIDs(ctx context.Context, ids []int) ([]PT, error) {
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

func (m *BaseManager[T, PT]) List(ctx context.Context, filter Filter) ([]PT, error) {
	entities, err := m.repo.List(ctx, filter)
	if err != nil {
		m.logger.Error("List failed", logger.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("list failed: %w", err)
	}
	return entities, nil
}

func (m *BaseManager[T, PT]) Count(ctx context.Context, filter Filter) (int, error) {
	count, err := m.repo.Count(ctx, filter)
	if err != nil {
		m.logger.Error("Count failed", logger.Error(err))
		return 0, fmt.Errorf("count failed: %w", err)
	}
	return count, nil
}

func (m *BaseManager[T, PT]) Exists(ctx context.Context, id int) (bool, error) {
	exists, err := m.repo.Exists(ctx, id)
	if err != nil {
		m.logger.Error("Exists check failed", logger.Error(err), logger.Int("id", id))
		return false, fmt.Errorf("exists check failed: %w", err)
	}
	return exists, nil
}

func (m *BaseManager[T, PT]) ExecuteInTx(
	ctx context.Context,
	txProvider TxProvider,
	ops func(repo Repository[T, PT]) error,
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

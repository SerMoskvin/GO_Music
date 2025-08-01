package db

import (
	"context"
	"database/sql"
)

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

// Repository интерфейс
type Repository[T any, ID comparable] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ID) error
	GetByID(ctx context.Context, id ID) (*T, error)
	GetByIDs(ctx context.Context, ids []ID) ([]*T, error)
	List(ctx context.Context, filter Filter) ([]*T, error)
	Count(ctx context.Context, filter Filter) (int, error)
	Exists(ctx context.Context, id ID) (bool, error)
	WithTx(tx *sql.Tx) Repository[T, ID]
}

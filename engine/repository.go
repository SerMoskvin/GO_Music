package engine

import (
	"errors"
)

// Entity определяет базовые методы для всех сущностей
type Entity interface {
	GetID() int
	SetID(id int)
	Validate() error
}

// PointerEntity - вспомогательный интерфейс для работы с указателями
type PointerEntity[T any] interface {
	*T
	Entity
}

// Repository - обобщенный интерфейс репозитория
type Repository[T any, PT PointerEntity[T]] interface {
	Create(entity PT) error
	Update(entity PT) error
	Delete(id int) error
	GetByID(id int) (PT, error)
	GetByIDs(ids []int) ([]PT, error)
}

// BaseManager - базовая реализация CRUD операций
type BaseManager[T any, PT PointerEntity[T]] struct {
	repo Repository[T, PT]
}

// NewBaseManager создает новый экземпляр BaseManager
func NewBaseManager[T any, PT PointerEntity[T]](repo Repository[T, PT]) *BaseManager[T, PT] {
	return &BaseManager[T, PT]{repo: repo}
}

// Create реализует создание сущности
func (m *BaseManager[T, PT]) Create(entity PT) error {
	if err := entity.Validate(); err != nil {
		return err
	}
	return m.repo.Create(entity)
}

// Update реализует обновление сущности
func (m *BaseManager[T, PT]) Update(entity PT) error {
	if entity.GetID() == 0 {
		return errors.New("ID не указан")
	}
	if err := entity.Validate(); err != nil {
		return err
	}
	return m.repo.Update(entity)
}

// Delete реализует удаление сущности
func (m *BaseManager[T, PT]) Delete(id int) error {
	if id == 0 {
		return errors.New("ID не указан")
	}
	return m.repo.Delete(id)
}

// GetByID реализует получение сущности по ID
func (m *BaseManager[T, PT]) GetByID(id int) (PT, error) {
	if id == 0 {
		return nil, errors.New("ID не указан")
	}
	return m.repo.GetByID(id)
}

// GetByIDs реализует получение нескольких сущностей по IDs
func (m *BaseManager[T, PT]) GetByIDs(ids []int) ([]PT, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

package engine

type Repository[T any] interface {
	Create(entity *T) error
	Update(entity *T) error
	Delete(id int) error
	GetByID(id int) (*T, error)
	GetByIDs(ids []int) ([]*T, error)
}

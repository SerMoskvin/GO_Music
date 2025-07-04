package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type UserRepository interface {
	Repository[domain.User]
}

type UserManager struct {
	repo UserRepository
}

func NewUserManager(repo UserRepository) *UserManager {
	return &UserManager{repo: repo}
}

func (m *UserManager) Create(user *domain.User) error {
	if user == nil {
		return errors.New("пользователь не указан")
	}
	if err := validate.ValidateStruct(user); err != nil {
		return err
	}
	return m.repo.Create(user)
}

func (m *UserManager) Update(user *domain.User) error {
	if user == nil {
		return errors.New("пользователь не указан")
	}
	if user.UserID == 0 {
		return errors.New("не указан ID пользователя")
	}
	if err := validate.ValidateStruct(user); err != nil {
		return err
	}
	return m.repo.Update(user)
}

func (m *UserManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID пользователя")
	}
	return m.repo.Delete(id)
}

func (m *UserManager) GetByID(id int) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("не указан ID пользователя")
	}
	return m.repo.GetByID(id)
}

func (m *UserManager) GetByIDs(ids []int) ([]*domain.User, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

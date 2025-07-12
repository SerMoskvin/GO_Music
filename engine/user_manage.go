package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/access"

	"github.com/SerMoskvin/validate"
)

type UserRepository interface {
	Repository[domain.User]
	GetByLogin(login string) (*domain.User, error)
}

type UserManager struct {
	repo UserRepository
	auth *access.Authenticator
}

func NewUserManager(repo UserRepository, auth *access.Authenticator) *UserManager {
	return &UserManager{
		repo: repo,
		auth: auth,
	}
}

func (m *UserManager) Create(user *domain.User) error {
	if user == nil {
		return errors.New("пользователь не указан")
	}
	if err := validate.ValidateStruct(user); err != nil {
		return err
	}
	if user.Password != "" {
		hashed, err := m.auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
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
	if user.Password != "" {
		hashed, err := m.auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
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

func (m *UserManager) Login(login, password string) (string, error) {
	if login == "" || password == "" {
		return "", errors.New("логин и пароль обязательны")
	}

	user, err := m.repo.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("пользователь не найден")
	}

	if !m.auth.CheckPasswordHash(password, user.Password) {
		return "", errors.New("неверный пароль")
	}

	token, err := m.auth.GenerateJWT(user.UserID, user.Login, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

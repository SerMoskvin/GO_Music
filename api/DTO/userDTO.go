package dto

import (
	"GO_Music/domain"
	"GO_Music/engine"
)

// UserCreateDTO для создания пользователя
type UserCreateDTO struct {
	Login    string `json:"login" validate:"required,min=1,max=250"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,min=1,max=50"`
	Surname  string `json:"surname" validate:"required,min=1,max=100"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Image    []byte `json:"image,omitempty"`
}

// UserUpdateDTO для обновления пользователя
type UserUpdateDTO struct {
	Login    *string `json:"login,omitempty" validate:"omitempty,min=1,max=250"`
	Password *string `json:"password,omitempty" validate:"omitempty"`
	Role     *string `json:"role,omitempty" validate:"omitempty,min=1,max=50"`
	Surname  *string `json:"surname,omitempty" validate:"omitempty,min=1,max=100"`
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Image    []byte  `json:"image,omitempty"`
}

// UserResponseDTO для ответа API
type UserResponseDTO struct {
	UserID           int    `json:"user_id"`
	Login            string `json:"login"`
	Role             string `json:"role"`
	Surname          string `json:"surname"`
	Name             string `json:"name"`
	RegistrationDate string `json:"registration_date"`
	Email            string `json:"email"`
	Image            []byte `json:"image,omitempty"`
}

// UserLoginDTO для аутентификации
type UserLoginDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserChangePasswordDTO для смены пароля
type UserChangePasswordDTO struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

// UserMapper реализует маппинг для пользователей
type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDomain(dto *UserCreateDTO) *domain.User {
	return &domain.User{
		Login:    dto.Login,
		Password: dto.Password,
		Role:     dto.Role,
		Surname:  dto.Surname,
		Name:     dto.Name,
		Email:    dto.Email,
		Image:    dto.Image,
	}
}

func (m *UserMapper) UpdateDomain(user *domain.User, dto *UserUpdateDTO) {
	if dto.Login != nil {
		user.Login = *dto.Login
	}
	if dto.Password != nil {
		user.Password = *dto.Password
	}
	if dto.Role != nil {
		user.Role = *dto.Role
	}
	if dto.Surname != nil {
		user.Surname = *dto.Surname
	}
	if dto.Name != nil {
		user.Name = *dto.Name
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Image != nil {
		user.Image = dto.Image
	}
}

func (m *UserMapper) ToResponse(user *domain.User) *UserResponseDTO {
	// Убедимся, что есть изображение
	image := user.Image
	if len(image) == 0 {
		image = engine.DefaultImage
	}

	return &UserResponseDTO{
		UserID:           user.UserID,
		Login:            user.Login,
		Role:             user.Role,
		Surname:          user.Surname,
		Name:             user.Name,
		RegistrationDate: domain.ToDateTime(user.RegistrationDate),
		Email:            user.Email,
		Image:            image,
	}
}

func (m *UserMapper) ToResponseList(users []*domain.User) []*UserResponseDTO {
	result := make([]*UserResponseDTO, len(users))
	for i, user := range users {
		result[i] = m.ToResponse(user)
	}
	return result
}

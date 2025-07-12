package domain

import (
	"time"

	"github.com/SerMoskvin/validate"
)

// User представляет запись пользователя
type User struct {
	UserID           int       `json:"user_id"`
	Login            string    `json:"login" validate:"required,min=1,max=250"`
	Password         string    `json:"password" validate:"required"`
	Role             string    `json:"role" validate:"required,min=1,max=50"`
	Surname          string    `json:"surname" validate:"required,min=1,max=100"`
	Name             string    `json:"name" validate:"required,min=1,max=100"`
	RegistrationDate time.Time `json:"registration_date" validate:"required"`
	Image            []byte    `json:"image,omitempty"`
}

func (u *User) GetID() int {
	return u.UserID
}

func (u *User) SetID(id int) {
	u.UserID = id
}

func (u *User) Validate() error {
	return validate.ValidateStruct(u)
}

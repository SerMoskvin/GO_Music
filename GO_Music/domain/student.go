package domain

import (
	"time"

	"github.com/SerMoskvin/validate"
)

type Student struct {
	StudentID     int       `json:"student_id"`
	UserID        *int      `json:"user_id,omitempty"`
	Surname       string    `json:"surname" validate:"required,min=1,max=60"`
	Name          string    `json:"name" validate:"required,min=1,max=45"`
	FatherName    *string   `json:"father_name,omitempty" validate:"omitempty,max=55"`
	Birthday      time.Time `json:"birthday" validate:"required,birthday_past"`
	PhoneNumber   *string   `json:"phone_number,omitempty" validate:"omitempty,len=11"`
	GroupID       int       `json:"group_id" validate:"required"`
	MusprogrammID int       `json:"musprogramm_id" validate:"required"`
}

func (s *Student) GetID() int {
	return s.StudentID
}

func (s *Student) SetID(id int) {
	s.StudentID = id
}

func (s *Student) Validate() error {
	return validate.ValidateStruct(s)
}

package domain

import (
	"time"

	"github.com/SerMoskvin/validate"
)

// Employee представляет запись сотрудника из БД
type Employee struct {
	EmployeeID     int       `json:"employee_id"`
	UserID         *int      `json:"user_id,omitempty"`
	Surname        string    `json:"surname" validate:"required,min=1,max=60"`
	Name           string    `json:"name" validate:"required,min=1,max=45"`
	FatherName     *string   `json:"father_name,omitempty" validate:"omitempty,max=55"`
	Birthday       time.Time `json:"birthday" validate:"required,birthday_past"`
	PhoneNumber    string    `json:"phone_number" validate:"required,len=11"`
	Job            string    `json:"job" validate:"required,min=1,max=60"`
	WorkExperience int       `json:"work_experience" validate:"required,gte=0"`
}

// GetID возвращает идентификатор сотрудника
func (e *Employee) GetID() int {
	return e.EmployeeID
}

// SetID устанавливает идентификатор сотрудника
func (e *Employee) SetID(id int) {
	e.EmployeeID = id
}

// Validate запускает валидацию полей сотрудника
func (e *Employee) Validate() error {
	return validate.ValidateStruct(e)
}

package dto

import (
	"GO_Music/domain"
	"time"
)

// EmployeeCreateDTO для создания сотрудника
type EmployeeCreateDTO struct {
	UserID         *int      `json:"user_id,omitempty" validate:"omitempty"`
	Surname        string    `json:"surname" validate:"required,min=1,max=60"`
	Name           string    `json:"name" validate:"required,min=1,max=45"`
	FatherName     *string   `json:"father_name,omitempty" validate:"omitempty,max=55"`
	Birthday       time.Time `json:"birthday" validate:"required,birthday_past"`
	PhoneNumber    string    `json:"phone_number" validate:"required,len=11"`
	Job            string    `json:"job" validate:"required,min=1,max=60"`
	WorkExperience int       `json:"work_experience" validate:"required,gte=0"`
}

// EmployeeUpdateDTO для обновления сотрудника
type EmployeeUpdateDTO struct {
	UserID         *int       `json:"user_id,omitempty" validate:"omitempty"`
	Surname        *string    `json:"surname,omitempty" validate:"omitempty,min=1,max=60"`
	Name           *string    `json:"name,omitempty" validate:"omitempty,min=1,max=45"`
	FatherName     *string    `json:"father_name,omitempty" validate:"omitempty,max=55"`
	Birthday       *time.Time `json:"birthday,omitempty" validate:"omitempty,birthday_past"`
	PhoneNumber    *string    `json:"phone_number,omitempty" validate:"omitempty,len=11"`
	Job            *string    `json:"job,omitempty" validate:"omitempty,min=1,max=60"`
	WorkExperience *int       `json:"work_experience,omitempty" validate:"omitempty,gte=0"`
}

// EmployeeResponseDTO для ответа API
type EmployeeResponseDTO struct {
	EmployeeID     int     `json:"employee_id"`
	UserID         *int    `json:"user_id,omitempty"`
	Surname        string  `json:"surname"`
	Name           string  `json:"name"`
	FatherName     *string `json:"father_name,omitempty"`
	Birthday       string  `json:"birthday"`
	PhoneNumber    string  `json:"phone_number"`
	Job            string  `json:"job"`
	WorkExperience int     `json:"work_experience"`
}

// EmployeeMapper реализует маппинг для сотрудников
type EmployeeMapper struct{}

func NewEmployeeMapper() *EmployeeMapper {
	return &EmployeeMapper{}
}

func (m *EmployeeMapper) ToDomain(dto *EmployeeCreateDTO) *domain.Employee {
	return &domain.Employee{
		UserID:         dto.UserID,
		Surname:        dto.Surname,
		Name:           dto.Name,
		FatherName:     dto.FatherName,
		Birthday:       dto.Birthday,
		PhoneNumber:    dto.PhoneNumber,
		Job:            dto.Job,
		WorkExperience: dto.WorkExperience,
	}
}

func (m *EmployeeMapper) UpdateDomain(employee *domain.Employee, dto *EmployeeUpdateDTO) {
	if dto.UserID != nil {
		employee.UserID = dto.UserID
	}
	if dto.Surname != nil {
		employee.Surname = *dto.Surname
	}
	if dto.Name != nil {
		employee.Name = *dto.Name
	}
	if dto.FatherName != nil {
		employee.FatherName = dto.FatherName
	}
	if dto.Birthday != nil {
		employee.Birthday = *dto.Birthday
	}
	if dto.PhoneNumber != nil {
		employee.PhoneNumber = *dto.PhoneNumber
	}
	if dto.Job != nil {
		employee.Job = *dto.Job
	}
	if dto.WorkExperience != nil {
		employee.WorkExperience = *dto.WorkExperience
	}
}

func (m *EmployeeMapper) ToResponse(employee *domain.Employee) *EmployeeResponseDTO {
	return &EmployeeResponseDTO{
		EmployeeID:     employee.EmployeeID,
		UserID:         employee.UserID,
		Surname:        employee.Surname,
		Name:           employee.Name,
		FatherName:     employee.FatherName,
		Birthday:       employee.Birthday.Format("2006-01-02"),
		PhoneNumber:    employee.PhoneNumber,
		Job:            employee.Job,
		WorkExperience: employee.WorkExperience,
	}
}

func (m *EmployeeMapper) ToResponseList(employees []*domain.Employee) []*EmployeeResponseDTO {
	result := make([]*EmployeeResponseDTO, len(employees))
	for i, emp := range employees {
		result[i] = m.ToResponse(emp)
	}
	return result
}

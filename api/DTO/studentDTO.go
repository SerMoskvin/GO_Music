package dto

import (
	"GO_Music/domain"
	"time"
)

// StudentCreateDTO для создания студента
type StudentCreateDTO struct {
	UserID        *int      `json:"user_id,omitempty" validate:"omitempty"`
	Surname       string    `json:"surname" validate:"required,min=1,max=60"`
	Name          string    `json:"name" validate:"required,min=1,max=45"`
	FatherName    *string   `json:"father_name,omitempty" validate:"omitempty,max=55"`
	Birthday      time.Time `json:"birthday" validate:"required,birthday_past"`
	PhoneNumber   *string   `json:"phone_number,omitempty" validate:"omitempty,len=11"`
	GroupID       int       `json:"group_id" validate:"required"`
	MusprogrammID int       `json:"musprogramm_id" validate:"required"`
}

// StudentUpdateDTO для обновления студента
type StudentUpdateDTO struct {
	UserID        *int       `json:"user_id,omitempty" validate:"omitempty"`
	Surname       *string    `json:"surname,omitempty" validate:"omitempty,min=1,max=60"`
	Name          *string    `json:"name,omitempty" validate:"omitempty,min=1,max=45"`
	FatherName    *string    `json:"father_name,omitempty" validate:"omitempty,max=55"`
	Birthday      *time.Time `json:"birthday,omitempty" validate:"omitempty,birthday_past"`
	PhoneNumber   *string    `json:"phone_number,omitempty" validate:"omitempty,len=11"`
	GroupID       *int       `json:"group_id,omitempty" validate:"omitempty"`
	MusprogrammID *int       `json:"musprogramm_id,omitempty" validate:"omitempty"`
}

// StudentResponseDTO для ответа API
type StudentResponseDTO struct {
	StudentID     int     `json:"student_id"`
	UserID        *int    `json:"user_id,omitempty"`
	Surname       string  `json:"surname"`
	Name          string  `json:"name"`
	FatherName    *string `json:"father_name,omitempty"`
	Birthday      string  `json:"birthday"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	GroupID       int     `json:"group_id"`
	MusprogrammID int     `json:"musprogramm_id"`
}

// StudentMapper реализует маппинг для студентов
type StudentMapper struct{}

func NewStudentMapper() *StudentMapper {
	return &StudentMapper{}
}

func (m *StudentMapper) ToDomain(dto *StudentCreateDTO) *domain.Student {
	return &domain.Student{
		UserID:        dto.UserID,
		Surname:       dto.Surname,
		Name:          dto.Name,
		FatherName:    dto.FatherName,
		Birthday:      dto.Birthday,
		PhoneNumber:   dto.PhoneNumber,
		GroupID:       dto.GroupID,
		MusprogrammID: dto.MusprogrammID,
	}
}

func (m *StudentMapper) UpdateDomain(student *domain.Student, dto *StudentUpdateDTO) {
	if dto.UserID != nil {
		student.UserID = dto.UserID
	}
	if dto.Surname != nil {
		student.Surname = *dto.Surname
	}
	if dto.Name != nil {
		student.Name = *dto.Name
	}
	if dto.FatherName != nil {
		student.FatherName = dto.FatherName
	}
	if dto.Birthday != nil {
		student.Birthday = *dto.Birthday
	}
	if dto.PhoneNumber != nil {
		student.PhoneNumber = dto.PhoneNumber
	}
	if dto.GroupID != nil {
		student.GroupID = *dto.GroupID
	}
	if dto.MusprogrammID != nil {
		student.MusprogrammID = *dto.MusprogrammID
	}
}

func (m *StudentMapper) ToResponse(student *domain.Student) *StudentResponseDTO {
	return &StudentResponseDTO{
		StudentID:     student.StudentID,
		UserID:        student.UserID,
		Surname:       student.Surname,
		Name:          student.Name,
		FatherName:    student.FatherName,
		Birthday:      domain.ToDMY(student.Birthday),
		PhoneNumber:   student.PhoneNumber,
		GroupID:       student.GroupID,
		MusprogrammID: student.MusprogrammID,
	}
}

func (m *StudentMapper) ToResponseList(students []*domain.Student) []*StudentResponseDTO {
	result := make([]*StudentResponseDTO, len(students))
	for i, student := range students {
		result[i] = m.ToResponse(student)
	}
	return result
}

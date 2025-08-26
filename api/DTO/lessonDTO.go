package dto

import (
	"GO_Music/domain"
)

// LessonCreateDTO для создания занятия
type LessonCreateDTO struct {
	AudienceID *int   `json:"audience_id,omitempty" validate:"omitempty"`
	EmployeeID int    `json:"employee_id" validate:"required"`
	GroupID    int    `json:"group_id" validate:"required"`
	StudentID  *int   `json:"student_id,omitempty" validate:"omitempty"`
	LessonName string `json:"lesson_name" validate:"required,min=1,max=70"`
	SubjectID  int    `json:"subject_id" validate:"required"`
}

// LessonUpdateDTO для обновления занятия
type LessonUpdateDTO struct {
	AudienceID *int    `json:"audience_id,omitempty" validate:"omitempty"`
	EmployeeID *int    `json:"employee_id,omitempty" validate:"omitempty"`
	GroupID    *int    `json:"group_id,omitempty" validate:"omitempty"`
	StudentID  *int    `json:"student_id,omitempty" validate:"omitempty"`
	LessonName *string `json:"lesson_name,omitempty" validate:"omitempty,min=1,max=70"`
	SubjectID  *int    `json:"subject_id,omitempty" validate:"omitempty"`
}

// LessonResponseDTO для ответа API
type LessonResponseDTO struct {
	LessonID   int    `json:"lesson_id"`
	AudienceID *int   `json:"audience_id,omitempty"`
	EmployeeID int    `json:"employee_id"`
	GroupID    int    `json:"group_id"`
	StudentID  *int   `json:"student_id,omitempty"`
	LessonName string `json:"lesson_name"`
	SubjectID  int    `json:"subject_id"`
}

// LessonMapper реализует маппинг для занятий
type LessonMapper struct{}

func NewLessonMapper() *LessonMapper {
	return &LessonMapper{}
}

func (m *LessonMapper) ToDomain(dto *LessonCreateDTO) *domain.Lesson {
	return &domain.Lesson{
		AudienceID: dto.AudienceID,
		EmployeeID: dto.EmployeeID,
		GroupID:    dto.GroupID,
		StudentID:  dto.StudentID,
		LessonName: dto.LessonName,
		SubjectID:  dto.SubjectID,
	}
}

func (m *LessonMapper) UpdateDomain(lesson *domain.Lesson, dto *LessonUpdateDTO) {
	if dto.AudienceID != nil {
		lesson.AudienceID = dto.AudienceID
	}
	if dto.EmployeeID != nil {
		lesson.EmployeeID = *dto.EmployeeID
	}
	if dto.GroupID != nil {
		lesson.GroupID = *dto.GroupID
	}
	if dto.StudentID != nil {
		lesson.StudentID = dto.StudentID
	}
	if dto.LessonName != nil {
		lesson.LessonName = *dto.LessonName
	}
	if dto.SubjectID != nil {
		lesson.SubjectID = *dto.SubjectID
	}
}

func (m *LessonMapper) ToResponse(lesson *domain.Lesson) *LessonResponseDTO {
	return &LessonResponseDTO{
		LessonID:   lesson.LessonID,
		AudienceID: lesson.AudienceID,
		EmployeeID: lesson.EmployeeID,
		GroupID:    lesson.GroupID,
		StudentID:  lesson.StudentID,
		LessonName: lesson.LessonName,
		SubjectID:  lesson.SubjectID,
	}
}

func (m *LessonMapper) ToResponseList(lessons []*domain.Lesson) []*LessonResponseDTO {
	result := make([]*LessonResponseDTO, len(lessons))
	for i, lesson := range lessons {
		result[i] = m.ToResponse(lesson)
	}
	return result
}

package dto

import (
	"GO_Music/domain"
	"time"
)

type StudentAssessmentCreateDTO struct {
	LessonID       int    `json:"lesson_id" validate:"required"`
	StudentID      int    `json:"student_id" validate:"required"`
	TaskType       string `json:"task_type" validate:"required,min=1,max=70"`
	Grade          int    `json:"grade" validate:"required"`
	AssessmentDate string `json:"assessment_date" validate:"required"` // Строка в формате DD.MM.YYYY
}

type StudentAssessmentUpdateDTO struct {
	LessonID       *int    `json:"lesson_id,omitempty" validate:"omitempty"`
	StudentID      *int    `json:"student_id,omitempty" validate:"omitempty"`
	TaskType       *string `json:"task_type,omitempty" validate:"omitempty,min=1,max=70"`
	Grade          *int    `json:"grade,omitempty" validate:"omitempty"`
	AssessmentDate *string `json:"assessment_date,omitempty" validate:"omitempty"` // Строка в формате DD.MM.YYYY
}

type StudentAssessmentResponseDTO struct {
	ID             int    `json:"id"`
	LessonID       int    `json:"lesson_id"`
	StudentID      int    `json:"student_id"`
	TaskType       string `json:"task_type"`
	Grade          int    `json:"grade"`
	AssessmentDate string `json:"assessment_date"` // Строка в формате DD.MM.YYYY
}

// AssessmentMapper реализует маппинг для оценок
type AssessmentMapper struct{}

// NewAssessmentMapper создает новый маппер для оценок
func NewAssessmentMapper() *AssessmentMapper {
	return &AssessmentMapper{}
}

// ParseDMY преобразует "DD.MM.YYYY" в time.Time
func ParseDMY(dateStr string) time.Time {
	t, _ := time.Parse("02.01.2006", dateStr)
	return t
}

// ToDMY преобразует time.Time в "DD.MM.YYYY"
func ToDMY(t time.Time) string {
	return t.Format("02.01.2006")
}

// ToDomain преобразует CreateDTO в доменную модель
func (m *AssessmentMapper) ToDomain(dto *StudentAssessmentCreateDTO) *domain.StudentAssessment {
	return &domain.StudentAssessment{
		LessonID:       dto.LessonID,
		StudentID:      dto.StudentID,
		TaskType:       dto.TaskType,
		Grade:          dto.Grade,
		AssessmentDate: ParseDMY(dto.AssessmentDate), // Преобразуем строку в time.Time
	}
}

// UpdateDomain обновляет доменную модель из UpdateDTO
func (m *AssessmentMapper) UpdateDomain(assessment *domain.StudentAssessment, dto *StudentAssessmentUpdateDTO) {
	if dto.LessonID != nil {
		assessment.LessonID = *dto.LessonID
	}
	if dto.StudentID != nil {
		assessment.StudentID = *dto.StudentID
	}
	if dto.TaskType != nil {
		assessment.TaskType = *dto.TaskType
	}
	if dto.Grade != nil {
		assessment.Grade = *dto.Grade
	}
	if dto.AssessmentDate != nil {
		assessment.AssessmentDate = ParseDMY(*dto.AssessmentDate) // Преобразуем строку в time.Time
	}
}

// ToResponse преобразует доменную модель в ResponseDTO
func (m *AssessmentMapper) ToResponse(assessment *domain.StudentAssessment) *StudentAssessmentResponseDTO {
	return &StudentAssessmentResponseDTO{
		ID:             assessment.AssessmentNoteID,
		LessonID:       assessment.LessonID,
		StudentID:      assessment.StudentID,
		TaskType:       assessment.TaskType,
		Grade:          assessment.Grade,
		AssessmentDate: ToDMY(assessment.AssessmentDate), // Преобразуем time.Time в строку DD.MM.YYYY
	}
}

// ToResponseList преобразует список доменных моделей
func (m *AssessmentMapper) ToResponseList(assessments []*domain.StudentAssessment) []*StudentAssessmentResponseDTO {
	result := make([]*StudentAssessmentResponseDTO, len(assessments))
	for i, a := range assessments {
		result[i] = m.ToResponse(a)
	}
	return result
}

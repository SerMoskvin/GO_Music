package dto

import (
	"GO_Music/domain"
)

// SubjectCreateDTO для создания предмета
type SubjectCreateDTO struct {
	SubjectName string `json:"subject_name" validate:"required,min=1,max=60"`
	SubjectType string `json:"subject_type" validate:"required,min=1,max=30"`
	ShortDesc   string `json:"short_desc" validate:"required"`
}

// SubjectUpdateDTO для обновления предмета
type SubjectUpdateDTO struct {
	SubjectName *string `json:"subject_name,omitempty" validate:"omitempty,min=1,max=60"`
	SubjectType *string `json:"subject_type,omitempty" validate:"omitempty,min=1,max=30"`
	ShortDesc   *string `json:"short_desc,omitempty" validate:"omitempty"`
}

// SubjectResponseDTO для ответа API
type SubjectResponseDTO struct {
	SubjectID   int    `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	SubjectType string `json:"subject_type"`
	ShortDesc   string `json:"short_desc"`
}

// SubjectMapper реализует маппинг для предметов
type SubjectMapper struct{}

func NewSubjectMapper() *SubjectMapper {
	return &SubjectMapper{}
}

func (m *SubjectMapper) ToDomain(dto *SubjectCreateDTO) *domain.Subject {
	return &domain.Subject{
		SubjectName: dto.SubjectName,
		SubjectType: dto.SubjectType,
		ShortDesc:   dto.ShortDesc,
	}
}

func (m *SubjectMapper) UpdateDomain(subject *domain.Subject, dto *SubjectUpdateDTO) {
	if dto.SubjectName != nil {
		subject.SubjectName = *dto.SubjectName
	}
	if dto.SubjectType != nil {
		subject.SubjectType = *dto.SubjectType
	}
	if dto.ShortDesc != nil {
		subject.ShortDesc = *dto.ShortDesc
	}
}

func (m *SubjectMapper) ToResponse(subject *domain.Subject) *SubjectResponseDTO {
	return &SubjectResponseDTO{
		SubjectID:   subject.SubjectID,
		SubjectName: subject.SubjectName,
		SubjectType: subject.SubjectType,
		ShortDesc:   subject.ShortDesc,
	}
}

func (m *SubjectMapper) ToResponseList(subjects []*domain.Subject) []*SubjectResponseDTO {
	result := make([]*SubjectResponseDTO, len(subjects))
	for i, subject := range subjects {
		result[i] = m.ToResponse(subject)
	}
	return result
}

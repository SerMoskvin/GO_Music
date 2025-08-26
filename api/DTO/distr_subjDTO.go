package dto

import (
	"GO_Music/domain"
)

// SubjectDistributionCreateDTO для создания распределения предмета
type SubjectDistributionCreateDTO struct {
	EmployeeID int `json:"employee_id" validate:"required"`
	SubjectID  int `json:"subject_id" validate:"required"`
}

// SubjectDistributionUpdateDTO для обновления распределения предмета
type SubjectDistributionUpdateDTO struct {
	EmployeeID *int `json:"employee_id,omitempty" validate:"omitempty"`
	SubjectID  *int `json:"subject_id,omitempty" validate:"omitempty"`
}

// SubjectDistributionResponseDTO для ответа API
type SubjectDistributionResponseDTO struct {
	SubjectDistrID int `json:"subject_distr_id"`
	EmployeeID     int `json:"employee_id"`
	SubjectID      int `json:"subject_id"`
}

// SubjectDistributionMapper реализует маппинг для распределений предметов
type SubjectDistributionMapper struct{}

func NewSubjectDistributionMapper() *SubjectDistributionMapper {
	return &SubjectDistributionMapper{}
}

func (m *SubjectDistributionMapper) ToDomain(dto *SubjectDistributionCreateDTO) *domain.SubjectDistribution {
	return &domain.SubjectDistribution{
		EmployeeID: dto.EmployeeID,
		SubjectID:  dto.SubjectID,
	}
}

func (m *SubjectDistributionMapper) UpdateDomain(distribution *domain.SubjectDistribution, dto *SubjectDistributionUpdateDTO) {
	if dto.EmployeeID != nil {
		distribution.EmployeeID = *dto.EmployeeID
	}
	if dto.SubjectID != nil {
		distribution.SubjectID = *dto.SubjectID
	}
}

func (m *SubjectDistributionMapper) ToResponse(distribution *domain.SubjectDistribution) *SubjectDistributionResponseDTO {
	return &SubjectDistributionResponseDTO{
		SubjectDistrID: distribution.SubjectDistrID,
		EmployeeID:     distribution.EmployeeID,
		SubjectID:      distribution.SubjectID,
	}
}

func (m *SubjectDistributionMapper) ToResponseList(distributions []*domain.SubjectDistribution) []*SubjectDistributionResponseDTO {
	result := make([]*SubjectDistributionResponseDTO, len(distributions))
	for i, distr := range distributions {
		result[i] = m.ToResponse(distr)
	}
	return result
}

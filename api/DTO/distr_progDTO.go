package dto

import (
	"GO_Music/domain"
)

// ProgrammDistributionCreateDTO для создания распределения программы
type ProgrammDistributionCreateDTO struct {
	MusprogrammID int `json:"musprogramm_id" validate:"required"`
	SubjectID     int `json:"subject_id" validate:"required"`
}

// ProgrammDistributionUpdateDTO для обновления распределения программы
type ProgrammDistributionUpdateDTO struct {
	MusprogrammID *int `json:"musprogramm_id,omitempty" validate:"omitempty"`
	SubjectID     *int `json:"subject_id,omitempty" validate:"omitempty"`
}

// ProgrammDistributionResponseDTO для ответа API
type ProgrammDistributionResponseDTO struct {
	ProgrammDistrID int `json:"programm_distr_id"`
	MusprogrammID   int `json:"musprogramm_id"`
	SubjectID       int `json:"subject_id"`
}

// ProgrammDistributionMapper реализует маппинг для распределений программ
type ProgrammDistributionMapper struct{}

func NewProgrammDistributionMapper() *ProgrammDistributionMapper {
	return &ProgrammDistributionMapper{}
}

func (m *ProgrammDistributionMapper) ToDomain(dto *ProgrammDistributionCreateDTO) *domain.ProgrammDistribution {
	return &domain.ProgrammDistribution{
		MusprogrammID: dto.MusprogrammID,
		SubjectID:     dto.SubjectID,
	}
}

func (m *ProgrammDistributionMapper) UpdateDomain(distribution *domain.ProgrammDistribution, dto *ProgrammDistributionUpdateDTO) {
	if dto.MusprogrammID != nil {
		distribution.MusprogrammID = *dto.MusprogrammID
	}
	if dto.SubjectID != nil {
		distribution.SubjectID = *dto.SubjectID
	}
}

func (m *ProgrammDistributionMapper) ToResponse(distribution *domain.ProgrammDistribution) *ProgrammDistributionResponseDTO {
	return &ProgrammDistributionResponseDTO{
		ProgrammDistrID: distribution.ProgrammDistrID,
		MusprogrammID:   distribution.MusprogrammID,
		SubjectID:       distribution.SubjectID,
	}
}

func (m *ProgrammDistributionMapper) ToResponseList(distributions []*domain.ProgrammDistribution) []*ProgrammDistributionResponseDTO {
	result := make([]*ProgrammDistributionResponseDTO, len(distributions))
	for i, distr := range distributions {
		result[i] = m.ToResponse(distr)
	}
	return result
}

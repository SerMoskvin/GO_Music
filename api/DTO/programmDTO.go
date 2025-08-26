package dto

import (
	"GO_Music/domain"
)

// ProgrammCreateDTO для создания музыкальной программы
type ProgrammCreateDTO struct {
	ProgrammName           string  `json:"programm_name" validate:"required,min=1,max=100"`
	ProgrammType           string  `json:"programm_type" validate:"required,min=1,max=70"`
	Duration               int     `json:"duration" validate:"required,gte=0"`
	Instrument             *string `json:"instrument,omitempty" validate:"omitempty,max=100"`
	Description            *string `json:"description,omitempty"`
	StudyLoad              int     `json:"study_load" validate:"required,gte=0"`
	FinalCertificationForm string  `json:"final_certification_form" validate:"required,min=1,max=100"`
}

// ProgrammUpdateDTO для обновления музыкальной программы
type ProgrammUpdateDTO struct {
	ProgrammName           *string `json:"programm_name,omitempty" validate:"omitempty,min=1,max=100"`
	ProgrammType           *string `json:"programm_type,omitempty" validate:"omitempty,min=1,max=70"`
	Duration               *int    `json:"duration,omitempty" validate:"omitempty,gte=0"`
	Instrument             *string `json:"instrument,omitempty" validate:"omitempty,max=100"`
	Description            *string `json:"description,omitempty"`
	StudyLoad              *int    `json:"study_load,omitempty" validate:"omitempty,gte=0"`
	FinalCertificationForm *string `json:"final_certification_form,omitempty" validate:"omitempty,min=1,max=100"`
}

// ProgrammResponseDTO для ответа API
type ProgrammResponseDTO struct {
	MusprogrammID          int     `json:"musprogramm_id"`
	ProgrammName           string  `json:"programm_name"`
	ProgrammType           string  `json:"programm_type"`
	Duration               int     `json:"duration"`
	Instrument             *string `json:"instrument,omitempty"`
	Description            *string `json:"description,omitempty"`
	StudyLoad              int     `json:"study_load"`
	FinalCertificationForm string  `json:"final_certification_form"`
}

// ProgrammMapper реализует маппинг для музыкальных программ
type ProgrammMapper struct{}

func NewProgrammMapper() *ProgrammMapper {
	return &ProgrammMapper{}
}

func (m *ProgrammMapper) ToDomain(dto *ProgrammCreateDTO) *domain.Programm {
	return &domain.Programm{
		ProgrammName:           dto.ProgrammName,
		ProgrammType:           dto.ProgrammType,
		Duration:               dto.Duration,
		Instrument:             dto.Instrument,
		Description:            dto.Description,
		StudyLoad:              dto.StudyLoad,
		FinalCertificationForm: dto.FinalCertificationForm,
	}
}

func (m *ProgrammMapper) UpdateDomain(programm *domain.Programm, dto *ProgrammUpdateDTO) {
	if dto.ProgrammName != nil {
		programm.ProgrammName = *dto.ProgrammName
	}
	if dto.ProgrammType != nil {
		programm.ProgrammType = *dto.ProgrammType
	}
	if dto.Duration != nil {
		programm.Duration = *dto.Duration
	}
	if dto.Instrument != nil {
		programm.Instrument = dto.Instrument
	}
	if dto.Description != nil {
		programm.Description = dto.Description
	}
	if dto.StudyLoad != nil {
		programm.StudyLoad = *dto.StudyLoad
	}
	if dto.FinalCertificationForm != nil {
		programm.FinalCertificationForm = *dto.FinalCertificationForm
	}
}

func (m *ProgrammMapper) ToResponse(programm *domain.Programm) *ProgrammResponseDTO {
	return &ProgrammResponseDTO{
		MusprogrammID:          programm.MusprogrammID,
		ProgrammName:           programm.ProgrammName,
		ProgrammType:           programm.ProgrammType,
		Duration:               programm.Duration,
		Instrument:             programm.Instrument,
		Description:            programm.Description,
		StudyLoad:              programm.StudyLoad,
		FinalCertificationForm: programm.FinalCertificationForm,
	}
}

func (m *ProgrammMapper) ToResponseList(programms []*domain.Programm) []*ProgrammResponseDTO {
	result := make([]*ProgrammResponseDTO, len(programms))
	for i, prog := range programms {
		result[i] = m.ToResponse(prog)
	}
	return result
}

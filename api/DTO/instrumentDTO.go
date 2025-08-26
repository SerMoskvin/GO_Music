package dto

import (
	"GO_Music/domain"
)

// InstrumentCreateDTO для создания инструмента
type InstrumentCreateDTO struct {
	AudienceID int    `json:"audience_id" validate:"required"`
	Name       string `json:"name" validate:"required,min=1,max=150"`
	InstrType  string `json:"instr_type" validate:"required,min=1,max=70"`
	Condition  string `json:"condition" validate:"required,min=1,max=70"`
}

// InstrumentUpdateDTO для обновления инструмента
type InstrumentUpdateDTO struct {
	AudienceID *int    `json:"audience_id,omitempty" validate:"omitempty"`
	Name       *string `json:"name,omitempty" validate:"omitempty,min=1,max=150"`
	InstrType  *string `json:"instr_type,omitempty" validate:"omitempty,min=1,max=70"`
	Condition  *string `json:"condition,omitempty" validate:"omitempty,min=1,max=70"`
}

// InstrumentResponseDTO для ответа API
type InstrumentResponseDTO struct {
	InstrumentID int    `json:"instrument_id"`
	AudienceID   int    `json:"audience_id"`
	Name         string `json:"name"`
	InstrType    string `json:"instr_type"`
	Condition    string `json:"condition"`
}

// InstrumentMapper реализует маппинг для инструментов
type InstrumentMapper struct{}

func NewInstrumentMapper() *InstrumentMapper {
	return &InstrumentMapper{}
}

func (m *InstrumentMapper) ToDomain(dto *InstrumentCreateDTO) *domain.Instrument {
	return &domain.Instrument{
		AudienceID: dto.AudienceID,
		Name:       dto.Name,
		InstrType:  dto.InstrType,
		Condition:  dto.Condition,
	}
}

func (m *InstrumentMapper) UpdateDomain(instrument *domain.Instrument, dto *InstrumentUpdateDTO) {
	if dto.AudienceID != nil {
		instrument.AudienceID = *dto.AudienceID
	}
	if dto.Name != nil {
		instrument.Name = *dto.Name
	}
	if dto.InstrType != nil {
		instrument.InstrType = *dto.InstrType
	}
	if dto.Condition != nil {
		instrument.Condition = *dto.Condition
	}
}

func (m *InstrumentMapper) ToResponse(instrument *domain.Instrument) *InstrumentResponseDTO {
	return &InstrumentResponseDTO{
		InstrumentID: instrument.InstrumentID,
		AudienceID:   instrument.AudienceID,
		Name:         instrument.Name,
		InstrType:    instrument.InstrType,
		Condition:    instrument.Condition,
	}
}

func (m *InstrumentMapper) ToResponseList(instruments []*domain.Instrument) []*InstrumentResponseDTO {
	result := make([]*InstrumentResponseDTO, len(instruments))
	for i, instr := range instruments {
		result[i] = m.ToResponse(instr)
	}
	return result
}

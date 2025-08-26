package dto

import (
	"GO_Music/domain"
)

// AudienceCreateDTO для создания аудитории
type AudienceCreateDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	AudinType   string `json:"audin_type" validate:"required,min=1,max=50"`
	AudinNumber string `json:"audin_number" validate:"required,min=1,max=30"`
	Capacity    int    `json:"capacity" validate:"required,min=1"`
}

// AudienceUpdateDTO для обновления аудитории
type AudienceUpdateDTO struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=50"`
	AudinType   *string `json:"audin_type,omitempty" validate:"omitempty,min=1,max=50"`
	AudinNumber *string `json:"audin_number,omitempty" validate:"omitempty,min=1,max=30"`
	Capacity    *int    `json:"capacity,omitempty" validate:"omitempty,min=1"`
}

// AudienceResponseDTO для ответа API
type AudienceResponseDTO struct {
	AudienceID  int    `json:"audience_id"`
	Name        string `json:"name"`
	AudinType   string `json:"audin_type"`
	AudinNumber string `json:"audin_number"`
	Capacity    int    `json:"capacity"`
}

// AudienceMapper реализует маппинг для аудиторий
type AudienceMapper struct{}

func NewAudienceMapper() *AudienceMapper {
	return &AudienceMapper{}
}

func (m *AudienceMapper) ToDomain(dto *AudienceCreateDTO) *domain.Audience {
	return &domain.Audience{
		Name:        dto.Name,
		AudinType:   dto.AudinType,
		AudinNumber: dto.AudinNumber,
		Capacity:    dto.Capacity,
	}
}

func (m *AudienceMapper) UpdateDomain(audience *domain.Audience, dto *AudienceUpdateDTO) {
	if dto.Name != nil {
		audience.Name = *dto.Name
	}
	if dto.AudinType != nil {
		audience.AudinType = *dto.AudinType
	}
	if dto.AudinNumber != nil {
		audience.AudinNumber = *dto.AudinNumber
	}
	if dto.Capacity != nil {
		audience.Capacity = *dto.Capacity
	}
}

func (m *AudienceMapper) ToResponse(audience *domain.Audience) *AudienceResponseDTO {
	return &AudienceResponseDTO{
		AudienceID:  audience.AudienceID,
		Name:        audience.Name,
		AudinType:   audience.AudinType,
		AudinNumber: audience.AudinNumber,
		Capacity:    audience.Capacity,
	}
}

func (m *AudienceMapper) ToResponseList(audiences []*domain.Audience) []*AudienceResponseDTO {
	result := make([]*AudienceResponseDTO, len(audiences))
	for i, audience := range audiences {
		result[i] = m.ToResponse(audience)
	}
	return result
}

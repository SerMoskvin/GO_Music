package dto

import (
	"GO_Music/domain"
)

// StudyGroupCreateDTO для создания учебной группы
type StudyGroupCreateDTO struct {
	MusProgrammID    int    `json:"musprogramm_id" validate:"required"`
	GroupName        string `json:"group_name" validate:"required,min=1,max=100"`
	StudyYear        int    `json:"study_year" validate:"required"`
	NumberOfStudents int    `json:"number_of_students" validate:"required"`
}

// StudyGroupUpdateDTO для обновления учебной группы
type StudyGroupUpdateDTO struct {
	MusProgrammID    *int    `json:"musprogramm_id,omitempty" validate:"omitempty"`
	GroupName        *string `json:"group_name,omitempty" validate:"omitempty,min=1,max=100"`
	StudyYear        *int    `json:"study_year,omitempty" validate:"omitempty"`
	NumberOfStudents *int    `json:"number_of_students,omitempty" validate:"omitempty"`
}

// StudyGroupResponseDTO для ответа API
type StudyGroupResponseDTO struct {
	GroupID          int    `json:"group_id"`
	MusProgrammID    int    `json:"musprogramm_id"`
	GroupName        string `json:"group_name"`
	StudyYear        int    `json:"study_year"`
	NumberOfStudents int    `json:"number_of_students"`
}

// StudyGroupMapper реализует маппинг для учебных групп
type StudyGroupMapper struct{}

func NewStudyGroupMapper() *StudyGroupMapper {
	return &StudyGroupMapper{}
}

func (m *StudyGroupMapper) ToDomain(dto *StudyGroupCreateDTO) *domain.StudyGroup {
	return &domain.StudyGroup{
		MusProgrammID:    dto.MusProgrammID,
		GroupName:        dto.GroupName,
		StudyYear:        dto.StudyYear,
		NumberOfStudents: dto.NumberOfStudents,
	}
}

func (m *StudyGroupMapper) UpdateDomain(group *domain.StudyGroup, dto *StudyGroupUpdateDTO) {
	if dto.MusProgrammID != nil {
		group.MusProgrammID = *dto.MusProgrammID
	}
	if dto.GroupName != nil {
		group.GroupName = *dto.GroupName
	}
	if dto.StudyYear != nil {
		group.StudyYear = *dto.StudyYear
	}
	if dto.NumberOfStudents != nil {
		group.NumberOfStudents = *dto.NumberOfStudents
	}
}

func (m *StudyGroupMapper) ToResponse(group *domain.StudyGroup) *StudyGroupResponseDTO {
	return &StudyGroupResponseDTO{
		GroupID:          group.GroupID,
		MusProgrammID:    group.MusProgrammID,
		GroupName:        group.GroupName,
		StudyYear:        group.StudyYear,
		NumberOfStudents: group.NumberOfStudents,
	}
}

func (m *StudyGroupMapper) ToResponseList(groups []*domain.StudyGroup) []*StudyGroupResponseDTO {
	result := make([]*StudyGroupResponseDTO, len(groups))
	for i, group := range groups {
		result[i] = m.ToResponse(group)
	}
	return result
}

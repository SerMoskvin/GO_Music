package dto

import (
	"GO_Music/domain"
)

// ScheduleCreateDTO DTO для создания расписания
type ScheduleCreateDTO struct {
	LessonID      int    `json:"lesson_id" validate:"required"`
	DayWeek       string `json:"day_week" validate:"required,min=1,max=20"`
	TimeBegin     string `json:"time_begin" validate:"required"`      // Формат "15:04"
	TimeEnd       string `json:"time_end" validate:"required"`        // Формат "15:04"
	SchdDateStart string `json:"schd_date_start" validate:"required"` // Формат "DD.MM.YYYY"
	SchdDateEnd   string `json:"schd_date_end" validate:"required"`   // Формат "DD.MM.YYYY"
}

// ScheduleUpdateDTO DTO для обновления расписания
type ScheduleUpdateDTO struct {
	LessonID      *int    `json:"lesson_id,omitempty"`
	DayWeek       *string `json:"day_week,omitempty"`
	TimeBegin     *string `json:"time_begin,omitempty"`      // Формат "15:04"
	TimeEnd       *string `json:"time_end,omitempty"`        // Формат "15:04"
	SchdDateStart *string `json:"schd_date_start,omitempty"` // Формат "DD.MM.YYYY"
	SchdDateEnd   *string `json:"schd_date_end,omitempty"`   // Формат "DD.MM.YYYY"
}

// ScheduleResponseDTO DTO для ответа с расписанием
type ScheduleResponseDTO struct {
	ScheduleID    int    `json:"schedule_id"`
	LessonID      int    `json:"lesson_id"`
	DayWeek       string `json:"day_week"`
	TimeBegin     string `json:"time_begin"`           // Формат "15:04"
	TimeEnd       string `json:"time_end"`             // Формат "15:04"
	SchdDateStart string `json:"schd_date_start"`      // Формат "DD.MM.YYYY"
	SchdDateEnd   string `json:"schd_date_end"`        // Формат "DD.MM.YYYY"
	CreatedAt     string `json:"created_at,omitempty"` // Формат "DD.MM.YYYY HH:MM:SS"
	UpdatedAt     string `json:"updated_at,omitempty"` // Формат "DD.MM.YYYY HH:MM:SS"
}

// ScheduleMapper маппер для расписания
type ScheduleMapper struct{}

// NewScheduleMapper создает новый маппер для расписания
func NewScheduleMapper() *ScheduleMapper {
	return &ScheduleMapper{}
}

// ToDomain преобразует CreateDTO в доменную модель
func (m *ScheduleMapper) ToDomain(dto *ScheduleCreateDTO) *domain.Schedule {
	return &domain.Schedule{
		LessonID:      dto.LessonID,
		DayWeek:       dto.DayWeek,
		TimeBegin:     domain.ParseTimeHM(dto.TimeBegin),
		TimeEnd:       domain.ParseTimeHM(dto.TimeEnd),
		SchdDateStart: domain.ParseDMY(dto.SchdDateStart),
		SchdDateEnd:   domain.ParseDMY(dto.SchdDateEnd),
	}
}

// UpdateDomain обновляет доменную модель из UpdateDTO
func (m *ScheduleMapper) UpdateDomain(schedule *domain.Schedule, dto *ScheduleUpdateDTO) {
	if dto.LessonID != nil {
		schedule.LessonID = *dto.LessonID
	}
	if dto.DayWeek != nil {
		schedule.DayWeek = *dto.DayWeek
	}
	if dto.TimeBegin != nil {
		schedule.TimeBegin = domain.ParseTimeHM(*dto.TimeBegin)
	}
	if dto.TimeEnd != nil {
		schedule.TimeEnd = domain.ParseTimeHM(*dto.TimeEnd)
	}
	if dto.SchdDateStart != nil {
		schedule.SchdDateStart = domain.ParseDMY(*dto.SchdDateStart)
	}
	if dto.SchdDateEnd != nil {
		schedule.SchdDateEnd = domain.ParseDMY(*dto.SchdDateEnd)
	}
}

// ToResponse преобразует доменную модель в ResponseDTO
func (m *ScheduleMapper) ToResponse(schedule *domain.Schedule) *ScheduleResponseDTO {
	return &ScheduleResponseDTO{
		ScheduleID:    schedule.ScheduleID,
		LessonID:      schedule.LessonID,
		DayWeek:       schedule.DayWeek,
		TimeBegin:     domain.ToTimeHM(schedule.TimeBegin),
		TimeEnd:       domain.ToTimeHM(schedule.TimeEnd),
		SchdDateStart: domain.ToDMY(schedule.SchdDateStart),
		SchdDateEnd:   domain.ToDMY(schedule.SchdDateEnd),
	}
}

// ToResponseList преобразует список доменных моделей в список ResponseDTO
func (m *ScheduleMapper) ToResponseList(schedules []*domain.Schedule) []*ScheduleResponseDTO {
	response := make([]*ScheduleResponseDTO, len(schedules))
	for i, schedule := range schedules {
		response[i] = m.ToResponse(schedule)
	}
	return response
}

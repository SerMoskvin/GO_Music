package domain

import (
	"time"

	"github.com/SerMoskvin/validate"
)

// Schedule представляет запись расписания
type Schedule struct {
	ScheduleID     int       `json:"schedule_id"`
	LessonID       int       `json:"lesson_id" validate:"required"`
	DayWeek        string    `json:"day_week" validate:"required,min=1,max=20"`
	TimeBegin      string    `json:"time_begin" validate:"required,datetime=15:04"`
	TimeEnd        string    `json:"time_end" validate:"required,datetime=15:04"`
	SchedDateStart time.Time `json:"sched_date_start" validate:"required"`
	SchedDateEnd   time.Time `json:"sched_date_end" validate:"required"`
}

func (s *Schedule) GetID() int {
	return s.ScheduleID
}

func (s *Schedule) SetID(id int) {
	s.ScheduleID = id
}

func (s *Schedule) Validate() error {
	return validate.ValidateStruct(s)
}

package domain

import (
	"github.com/SerMoskvin/validate"
)

// Lesson представляет запись занятия
type Lesson struct {
	LessonID   int    `json:"lesson_id"`
	AudienceID *int   `json:"audience_id,omitempty"`
	EmployeeID int    `json:"employee_id" validate:"required"`
	GroupID    int    `json:"group_id" validate:"required"`
	StudentID  *int   `json:"student_id,omitempty"`
	LessonName string `json:"lesson_name" validate:"required,min=1,max=70"`
	SubjectID  int    `json:"subject_id" validate:"required"`
}

func (l *Lesson) GetID() int {
	return l.LessonID
}

func (l *Lesson) SetID(id int) {
	l.LessonID = id
}

func (l *Lesson) Validate() error {
	return validate.ValidateStruct(l)
}

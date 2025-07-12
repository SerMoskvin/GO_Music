package domain

import (
	"github.com/SerMoskvin/validate"
)

// StudentAttendance представляет запись посещаемости студента
type StudentAttendance struct {
	AttendanceNoteID int    `json:"attendance_note_id"`
	StudentID        int    `json:"student_id" validate:"required"`
	LessonID         int    `json:"lesson_id" validate:"required"`
	PresenceMark     bool   `json:"presence_mark" validate:"required"`
	AttendanceDate   string `json:"attendance_date" validate:"required,datetime=2006-01-02"`
}

func (sa *StudentAttendance) GetID() int {
	return sa.AttendanceNoteID
}

func (sa *StudentAttendance) SetID(id int) {
	sa.AttendanceNoteID = id
}

func (sa *StudentAttendance) Validate() error {
	return validate.ValidateStruct(sa)
}

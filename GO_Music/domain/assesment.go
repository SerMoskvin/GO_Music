package domain

import (
	"github.com/SerMoskvin/validate"
)

// StudentAssessment представляет запись оценки студента
type StudentAssessment struct {
	AssessmentNoteID int    `json:"assessment_note_id"`
	LessonID         int    `json:"lesson_id" validate:"required"`
	StudentID        int    `json:"student_id" validate:"required"`
	TaskType         string `json:"task_type" validate:"required,min=1,max=70"`
	Grade            int    `json:"grade" validate:"required"`
	AssessmentDate   string `json:"assessment_date" validate:"required,datetime=2006-01-02"`
}

func (sa *StudentAssessment) GetID() int {
	return sa.AssessmentNoteID
}

func (sa *StudentAssessment) SetID(id int) {
	sa.AssessmentNoteID = id
}

func (sa *StudentAssessment) Validate() error {
	return validate.ValidateStruct(sa)
}

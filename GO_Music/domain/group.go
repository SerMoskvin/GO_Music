package domain

import (
	"github.com/SerMoskvin/validate"
)

// StudyGroup представляет запись группы
type StudyGroup struct {
	GroupID          int    `json:"group_id"`
	MusProgrammID    int    `json:"musprogramm_id" validate:"required"`
	GroupName        string `json:"group_name" validate:"required,min=1,max=100"`
	StudyYear        int    `json:"study_year" validate:"required"`
	NumberOfStudents int    `json:"number_of_students" validate:"required"`
}

func (g *StudyGroup) GetID() int {
	return g.GroupID
}

func (g *StudyGroup) SetID(id int) {
	g.GroupID = id
}

func (g *StudyGroup) Validate() error {
	return validate.ValidateStruct(g)
}

package domain

import (
	"github.com/SerMoskvin/validate"
)

// Programm представляет запись музыкальной программы
type Programm struct {
	MusprogrammID          int     `json:"musprogramm_id"`
	ProgrammName           string  `json:"programm_name" validate:"required,min=1,max=100"`
	ProgrammType           string  `json:"programm_type" validate:"required,min=1,max=70"`
	Duration               int     `json:"duration" validate:"required,gte=0"`
	Instrument             *string `json:"instrument,omitempty" validate:"omitempty,max=100"`
	Description            *string `json:"description,omitempty"`
	StudyLoad              int     `json:"study_load" validate:"required,gte=0"`
	FinalCertificationForm string  `json:"final_certification_form" validate:"required,min=1,max=100"`
}

func (p *Programm) GetID() int {
	return p.MusprogrammID
}

func (p *Programm) SetID(id int) {
	p.MusprogrammID = id
}

func (p *Programm) Validate() error {
	return validate.ValidateStruct(p)
}

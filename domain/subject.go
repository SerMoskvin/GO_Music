package domain

import (
	"github.com/SerMoskvin/validate"
)

// Subject представляет запись предмета
type Subject struct {
	SubjectID   int    `json:"subject_id"`
	SubjectName string `json:"subject_name" validate:"required,min=1,max=60"`
	SubjectType string `json:"subject_type" validate:"required,min=1,max=30"`
	ShortDesc   string `json:"short_desc" validate:"required"`
}

func (s *Subject) GetID() int {
	return s.SubjectID
}

func (s *Subject) SetID(id int) {
	s.SubjectID = id
}

func (s *Subject) Validate() error {
	return validate.ValidateStruct(s)
}

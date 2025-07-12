package domain

import (
	"github.com/SerMoskvin/validate"
)

// ProgrammDistribution представляет распределение предметов по программам
type ProgrammDistribution struct {
	ProgrammDistrID int `json:"programm_distr_id"`
	MusprogrammID   int `json:"musprogramm_id" validate:"required"`
	SubjectID       int `json:"subject_id" validate:"required"`
}

func (pd *ProgrammDistribution) GetID() int {
	return pd.ProgrammDistrID
}

func (pd *ProgrammDistribution) SetID(id int) {
	pd.ProgrammDistrID = id
}

func (pd *ProgrammDistribution) Validate() error {
	return validate.ValidateStruct(pd)
}

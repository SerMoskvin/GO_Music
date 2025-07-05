package domain

import (
	"github.com/SerMoskvin/validate"
)

// SubjectDistribution представляет распределение предметов по сотрудникам
type SubjectDistribution struct {
	SubjectDistrID int `json:"subject_distr_id"`
	EmployeeID     int `json:"employee_id" validate:"required"`
	SubjectID      int `json:"subject_id" validate:"required"`
}

func (sd *SubjectDistribution) GetID() int {
	return sd.SubjectDistrID
}

func (sd *SubjectDistribution) SetID(id int) {
	sd.SubjectDistrID = id
}

func (sd *SubjectDistribution) Validate() error {
	return validate.ValidateStruct(sd)
}

package domain

import (
	"github.com/SerMoskvin/validate"
)

// Audience представляет запись аудитории
type Audience struct {
	AudienceID  int    `json:"audience_id"`
	Name        string `json:"name" validate:"required,min=1,max=50"`
	AudinType   string `json:"audin_type" validate:"required,min=1,max=50"`
	AudinNumber string `json:"audin_number" validate:"required,min=1,max=30"`
	Capacity    int    `json:"capacity" validate:"required,min=1"`
}

func (a *Audience) GetID() int {
	return a.AudienceID
}

func (a *Audience) SetID(id int) {
	a.AudienceID = id
}

func (a *Audience) Validate() error {
	return validate.ValidateStruct(a)
}

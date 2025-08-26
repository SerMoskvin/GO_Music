package domain

import (
	"github.com/SerMoskvin/validate"
)

// Instrument представляет запись инструмента
type Instrument struct {
	InstrumentID int    `json:"instrument_id"`
	AudienceID   int    `json:"audience_id" validate:"required"`
	Name         string `json:"name" validate:"required,min=1,max=150"`
	InstrType    string `json:"instr_type" validate:"required,min=1,max=70"`
	Condition    string `json:"condition" validate:"required,min=1,max=70"`
}

func (i *Instrument) GetID() int {
	return i.InstrumentID
}

func (i *Instrument) SetID(id int) {
	i.InstrumentID = id
}

func (i *Instrument) Validate() error {
	return validate.ValidateStruct(i)
}

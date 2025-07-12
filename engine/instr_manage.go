package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type InstrumentRepository interface {
	Repository[domain.Instrument]
}

type InstrumentManager struct {
	repo InstrumentRepository
}

func NewInstrumentManager(repo InstrumentRepository) *InstrumentManager {
	return &InstrumentManager{repo: repo}
}

func (m *InstrumentManager) Create(instr *domain.Instrument) error {
	if instr == nil {
		return errors.New("инструмент не указан")
	}
	if err := validate.ValidateStruct(instr); err != nil {
		return err
	}
	return m.repo.Create(instr)
}

func (m *InstrumentManager) Update(instr *domain.Instrument) error {
	if instr == nil {
		return errors.New("инструмент не указан")
	}
	if instr.InstrumentID == 0 {
		return errors.New("не указан ID инструмента")
	}
	if err := validate.ValidateStruct(instr); err != nil {
		return err
	}
	return m.repo.Update(instr)
}

func (m *InstrumentManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID инструмента")
	}
	return m.repo.Delete(id)
}

func (m *InstrumentManager) GetByID(id int) (*domain.Instrument, error) {
	if id == 0 {
		return nil, errors.New("не указан ID инструмента")
	}
	return m.repo.GetByID(id)
}

func (m *InstrumentManager) GetByIDs(ids []int) ([]*domain.Instrument, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

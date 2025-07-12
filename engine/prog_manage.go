package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type ProgrammRepository interface {
	Repository[domain.Programm]
}

type ProgrammManager struct {
	repo ProgrammRepository
}

func NewProgrammManager(repo ProgrammRepository) *ProgrammManager {
	return &ProgrammManager{repo: repo}
}

func (m *ProgrammManager) Create(prog *domain.Programm) error {
	if prog == nil {
		return errors.New("programm is nil")
	}
	if err := validate.ValidateStruct(prog); err != nil {
		return err
	}
	return m.repo.Create(prog)
}

func (m *ProgrammManager) Update(prog *domain.Programm) error {
	if prog == nil {
		return errors.New("programm is nil")
	}
	if prog.MusprogrammID == 0 {
		return errors.New("не указан ID программы")
	}
	if err := validate.ValidateStruct(prog); err != nil {
		return err
	}
	return m.repo.Update(prog)
}

func (m *ProgrammManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID программы")
	}
	return m.repo.Delete(id)
}

func (m *ProgrammManager) GetByID(id int) (*domain.Programm, error) {
	if id == 0 {
		return nil, errors.New("не указан ID программы")
	}
	return m.repo.GetByID(id)
}

func (m *ProgrammManager) GetByIDs(ids []int) ([]*domain.Programm, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

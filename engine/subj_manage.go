package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type SubjectRepository interface {
	Repository[domain.Subject]
}

type SubjectManager struct {
	repo SubjectRepository
}

func NewSubjectManager(repo SubjectRepository) *SubjectManager {
	return &SubjectManager{repo: repo}
}

func (m *SubjectManager) Create(subject *domain.Subject) error {
	if err := validate.ValidateStruct(subject); err != nil {
		return err
	}
	return m.repo.Create(subject)
}

func (m *SubjectManager) Update(subject *domain.Subject) error {
	if subject.SubjectID == 0 {
		return errors.New("не указан ID предмета")
	}
	if err := validate.ValidateStruct(subject); err != nil {
		return err
	}
	return m.repo.Update(subject)
}

func (m *SubjectManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID предмета")
	}
	return m.repo.Delete(id)
}

func (m *SubjectManager) GetByID(id int) (*domain.Subject, error) {
	if id == 0 {
		return nil, errors.New("не указан ID предмета")
	}
	return m.repo.GetByID(id)
}

func (m *SubjectManager) GetByIDs(ids []int) ([]*domain.Subject, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type StudentAssessmentRepository interface {
	Repository[domain.StudentAssessment]
}

type StudentAssessmentManager struct {
	repo StudentAssessmentRepository
}

func NewStudentAssessmentManager(repo StudentAssessmentRepository) *StudentAssessmentManager {
	return &StudentAssessmentManager{repo: repo}
}

func (m *StudentAssessmentManager) Create(sa *domain.StudentAssessment) error {
	if sa == nil {
		return errors.New("student assessment is nil")
	}
	if err := validate.ValidateStruct(sa); err != nil {
		return err
	}
	return m.repo.Create(sa)
}

func (m *StudentAssessmentManager) Update(sa *domain.StudentAssessment) error {
	if sa == nil {
		return errors.New("student assessment is nil")
	}
	if sa.AssessmentNoteID == 0 {
		return errors.New("не указан ID оценки")
	}
	if err := validate.ValidateStruct(sa); err != nil {
		return err
	}
	return m.repo.Update(sa)
}

func (m *StudentAssessmentManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID оценки")
	}
	return m.repo.Delete(id)
}

func (m *StudentAssessmentManager) GetByID(id int) (*domain.StudentAssessment, error) {
	if id == 0 {
		return nil, errors.New("не указан ID оценки")
	}
	return m.repo.GetByID(id)
}

func (m *StudentAssessmentManager) GetByIDs(ids []int) ([]*domain.StudentAssessment, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

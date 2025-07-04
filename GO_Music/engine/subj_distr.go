package engine

import (
	"GO_Music/domain"
	"errors"
	"fmt"
)

type SubjectDistributionRepository interface {
	Create(sd *domain.SubjectDistribution) error
	Update(sd *domain.SubjectDistribution) error
	Delete(id int) error
	GetByID(id int) (*domain.SubjectDistribution, error)
	ExistsByEmployeeAndSubject(employeeID, subjectID int) (bool, error)
}

// Менеджер бизнес-логики для SubjectDistribution
type SubjectDistributionManager struct {
	repo SubjectDistributionRepository
}

func NewSubjectDistributionManager(repo SubjectDistributionRepository) *SubjectDistributionManager {
	return &SubjectDistributionManager{repo: repo}
}

func (m *SubjectDistributionManager) Create(sd *domain.SubjectDistribution) error {
	if sd == nil {
		return errors.New("subject distribution is nil")
	}
	if err := sd.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	exists, err := m.repo.ExistsByEmployeeAndSubject(sd.EmployeeID, sd.SubjectID)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return errors.New("this employee already has the specified subject assigned")
	}

	return m.repo.Create(sd)
}

func (m *SubjectDistributionManager) Update(sd *domain.SubjectDistribution) error {
	if sd == nil {
		return errors.New("subject distribution is nil")
	}
	if sd.SubjectDistrID == 0 {
		return errors.New("invalid subject_distr_id")
	}
	if err := sd.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	existing, err := m.repo.GetByID(sd.SubjectDistrID)
	if err != nil {
		return fmt.Errorf("failed to get existing record: %w", err)
	}
	if existing == nil {
		return errors.New("subject distribution not found")
	}

	exists, err := m.repo.ExistsByEmployeeAndSubject(sd.EmployeeID, sd.SubjectID)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists && !(existing.EmployeeID == sd.EmployeeID && existing.SubjectID == sd.SubjectID) {
		return errors.New("this employee already has the specified subject assigned")
	}

	return m.repo.Update(sd)
}

func (m *SubjectDistributionManager) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return m.repo.Delete(id)
}

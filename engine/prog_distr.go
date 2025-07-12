package engine

import (
	"GO_Music/domain"
	"errors"
	"fmt"
)

// Интерфейс репозитория для ProgrammDistribution
type ProgrammDistributionRepository interface {
	Create(pd *domain.ProgrammDistribution) error
	Update(pd *domain.ProgrammDistribution) error
	Delete(id int) error
	GetByID(id int) (*domain.ProgrammDistribution, error)
	ExistsByProgramAndSubject(musprogrammID, subjectID int) (bool, error)
}

// Менеджер бизнес-логики для ProgrammDistribution
type ProgrammDistributionManager struct {
	repo ProgrammDistributionRepository
}

func NewProgrammDistributionManager(repo ProgrammDistributionRepository) *ProgrammDistributionManager {
	return &ProgrammDistributionManager{repo: repo}
}

func (m *ProgrammDistributionManager) Create(pd *domain.ProgrammDistribution) error {
	if pd == nil {
		return errors.New("programm distribution is nil")
	}
	if err := pd.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	exists, err := m.repo.ExistsByProgramAndSubject(pd.MusprogrammID, pd.SubjectID)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return errors.New("this program already has the specified subject assigned")
	}

	return m.repo.Create(pd)
}

func (m *ProgrammDistributionManager) Update(pd *domain.ProgrammDistribution) error {
	if pd == nil {
		return errors.New("programm distribution is nil")
	}
	if pd.ProgrammDistrID == 0 {
		return errors.New("invalid programm_distr_id")
	}
	if err := pd.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	existing, err := m.repo.GetByID(pd.ProgrammDistrID)
	if err != nil {
		return fmt.Errorf("failed to get existing record: %w", err)
	}
	if existing == nil {
		return errors.New("programm distribution not found")
	}

	exists, err := m.repo.ExistsByProgramAndSubject(pd.MusprogrammID, pd.SubjectID)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists && !(existing.MusprogrammID == pd.MusprogrammID && existing.SubjectID == pd.SubjectID) {
		return errors.New("this program already has the specified subject assigned")
	}

	return m.repo.Update(pd)
}

func (m *ProgrammDistributionManager) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return m.repo.Delete(id)
}

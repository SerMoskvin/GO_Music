package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type AudienceRepository interface {
	Repository[domain.Audience]
}

type AudienceManager struct {
	repo AudienceRepository
}

func NewAudienceManager(repo AudienceRepository) *AudienceManager {
	return &AudienceManager{repo: repo}
}

func (m *AudienceManager) Create(audience *domain.Audience) error {
	if err := validate.ValidateStruct(audience); err != nil {
		return err
	}
	return m.repo.Create(audience)
}

func (m *AudienceManager) Update(audience *domain.Audience) error {
	if audience.AudienceID == 0 {
		return errors.New("не указан ID аудитории")
	}
	if err := validate.ValidateStruct(audience); err != nil {
		return err
	}
	return m.repo.Update(audience)
}

func (m *AudienceManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID аудитории")
	}
	return m.repo.Delete(id)
}

func (m *AudienceManager) GetByID(id int) (*domain.Audience, error) {
	if id == 0 {
		return nil, errors.New("не указан ID аудитории")
	}
	return m.repo.GetByID(id)
}

func (m *AudienceManager) GetByIDs(ids []int) ([]*domain.Audience, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

package engine

import (
	"GO_Music/domain"
	"errors"
)

type StudyGroupRepository interface {
	Repository[domain.StudyGroup]
}

type StudyGroupManager struct {
	repo StudyGroupRepository
}

func NewStudyGroupManager(repo StudyGroupRepository) *StudyGroupManager {
	return &StudyGroupManager{repo: repo}
}

func (m *StudyGroupManager) Create(group *domain.StudyGroup) error {
	if group == nil {
		return errors.New("группа обучения не указана")
	}
	if err := group.Validate(); err != nil {
		return err
	}
	return m.repo.Create(group)
}

func (m *StudyGroupManager) Update(group *domain.StudyGroup) error {
	if group == nil {
		return errors.New("группа обучения не указана")
	}
	if group.GroupID == 0 {
		return errors.New("не указан ID группы обучения")
	}
	if err := group.Validate(); err != nil {
		return err
	}
	return m.repo.Update(group)
}

func (m *StudyGroupManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID группы обучения")
	}
	return m.repo.Delete(id)
}

func (m *StudyGroupManager) GetByID(id int) (*domain.StudyGroup, error) {
	if id == 0 {
		return nil, errors.New("не указан ID группы обучения")
	}
	return m.repo.GetByID(id)
}

func (m *StudyGroupManager) GetByIDs(ids []int) ([]*domain.StudyGroup, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

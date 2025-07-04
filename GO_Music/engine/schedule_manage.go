package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type ScheduleRepository interface {
	Repository[domain.Schedule]
}

type ScheduleManager struct {
	repo ScheduleRepository
}

func NewScheduleManager(repo ScheduleRepository) *ScheduleManager {
	return &ScheduleManager{repo: repo}
}

func (m *ScheduleManager) Create(schedule *domain.Schedule) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}
	if err := validate.ValidateStruct(schedule); err != nil {
		return err
	}
	return m.repo.Create(schedule)
}

func (m *ScheduleManager) Update(schedule *domain.Schedule) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}
	if schedule.ScheduleID == 0 {
		return errors.New("не указан ID расписания")
	}
	if err := validate.ValidateStruct(schedule); err != nil {
		return err
	}
	return m.repo.Update(schedule)
}

func (m *ScheduleManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID расписания")
	}
	return m.repo.Delete(id)
}

func (m *ScheduleManager) GetByID(id int) (*domain.Schedule, error) {
	if id == 0 {
		return nil, errors.New("не указан ID расписания")
	}
	return m.repo.GetByID(id)
}

func (m *ScheduleManager) GetByIDs(ids []int) ([]*domain.Schedule, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

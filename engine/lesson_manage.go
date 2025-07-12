package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type LessonRepository interface {
	Repository[domain.Lesson]
}

type LessonManager struct {
	repo LessonRepository
}

func NewLessonManager(repo LessonRepository) *LessonManager {
	return &LessonManager{repo: repo}
}

func (m *LessonManager) Create(lesson *domain.Lesson) error {
	if lesson == nil {
		return errors.New("lesson is nil")
	}
	if err := validate.ValidateStruct(lesson); err != nil {
		return err
	}
	return m.repo.Create(lesson)
}

func (m *LessonManager) Update(lesson *domain.Lesson) error {
	if lesson == nil {
		return errors.New("lesson is nil")
	}
	if lesson.LessonID == 0 {
		return errors.New("не указан ID урока")
	}
	if err := validate.ValidateStruct(lesson); err != nil {
		return err
	}
	return m.repo.Update(lesson)
}

func (m *LessonManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID урока")
	}
	return m.repo.Delete(id)
}

func (m *LessonManager) GetByID(id int) (*domain.Lesson, error) {
	if id == 0 {
		return nil, errors.New("не указан ID урока")
	}
	return m.repo.GetByID(id)
}

func (m *LessonManager) GetByIDs(ids []int) ([]*domain.Lesson, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

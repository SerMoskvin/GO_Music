package engine

import (
	"GO_Music/domain"
	"errors"
)

type StudentRepository interface {
	Repository[domain.Student]
}

type StudentManager struct {
	repo      StudentRepository
	groupRepo StudyGroupRepository
}

func NewStudentManager(repo StudentRepository, groupRepo StudyGroupRepository) *StudentManager {
	return &StudentManager{
		repo:      repo,
		groupRepo: groupRepo,
	}
}

func (m *StudentManager) Create(student *domain.Student) error {
	if err := student.Validate(); err != nil {
		return err
	}
	return m.repo.Create(student)
}

func (m *StudentManager) Update(student *domain.Student) error {
	if student.StudentID == 0 {
		return errors.New("не указан ID студента")
	}
	if err := student.Validate(); err != nil {
		return err
	}
	return m.repo.Update(student)
}

func (m *StudentManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID студента")
	}
	return m.repo.Delete(id)
}

func (m *StudentManager) GetByID(id int) (*domain.Student, error) {
	if id == 0 {
		return nil, errors.New("не указан ID студента")
	}
	return m.repo.GetByID(id)
}

func (m *StudentManager) GetByIDs(ids []int) ([]*domain.Student, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

func (m *StudentManager) GetEducationalProgramByGroup(groupID int) (int, error) {
	group, err := m.groupRepo.GetByID(groupID)
	if err != nil {
		return 0, err
	}
	return group.MusProgrammID, nil
}

func (m *StudentManager) GetMusicalSpecialtyByStudent(studentID int) (int, error) {
	student, err := m.repo.GetByID(studentID)
	if err != nil {
		return 0, err
	}
	return student.MusprogrammID, nil
}

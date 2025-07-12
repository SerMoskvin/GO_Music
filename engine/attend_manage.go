package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type StudentAttendanceRepository interface {
	Repository[domain.StudentAttendance]
}

type StudentAttendanceManager struct {
	repo StudentAttendanceRepository
}

func NewStudentAttendanceManager(repo StudentAttendanceRepository) *StudentAttendanceManager {
	return &StudentAttendanceManager{repo: repo}
}

func (m *StudentAttendanceManager) Create(sa *domain.StudentAttendance) error {
	if sa == nil {
		return errors.New("student attendance is nil")
	}
	if err := validate.ValidateStruct(sa); err != nil {
		return err
	}
	return m.repo.Create(sa)
}

func (m *StudentAttendanceManager) Update(sa *domain.StudentAttendance) error {
	if sa == nil {
		return errors.New("student attendance is nil")
	}
	if sa.AttendanceNoteID == 0 {
		return errors.New("не указан ID записи посещения")
	}
	if err := validate.ValidateStruct(sa); err != nil {
		return err
	}
	return m.repo.Update(sa)
}

func (m *StudentAttendanceManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID записи посещения")
	}
	return m.repo.Delete(id)
}

func (m *StudentAttendanceManager) GetByID(id int) (*domain.StudentAttendance, error) {
	if id == 0 {
		return nil, errors.New("не указан ID записи посещения")
	}
	return m.repo.GetByID(id)
}

func (m *StudentAttendanceManager) GetByIDs(ids []int) ([]*domain.StudentAttendance, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

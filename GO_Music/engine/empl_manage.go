package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/validate"
)

type EmployeeRepository interface {
	Repository[domain.Employee]
}

type EmployeeManager struct {
	repo EmployeeRepository
}

func NewEmployeeManager(repo EmployeeRepository) *EmployeeManager {
	return &EmployeeManager{repo: repo}
}

func (m *EmployeeManager) Create(employee *domain.Employee) error {
	if employee == nil {
		return errors.New("employee is nil")
	}
	if err := validate.ValidateStruct(employee); err != nil {
		return err
	}
	return m.repo.Create(employee)
}

func (m *EmployeeManager) Update(employee *domain.Employee) error {
	if employee == nil {
		return errors.New("employee is nil")
	}
	if employee.EmployeeID == 0 {
		return errors.New("не указан ID сотрудника")
	}
	if err := validate.ValidateStruct(employee); err != nil {
		return err
	}
	return m.repo.Update(employee)
}

func (m *EmployeeManager) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID сотрудника")
	}
	return m.repo.Delete(id)
}

func (m *EmployeeManager) GetByID(id int) (*domain.Employee, error) {
	if id == 0 {
		return nil, errors.New("не указан ID сотрудника")
	}
	return m.repo.GetByID(id)
}

func (m *EmployeeManager) GetByIDs(ids []int) ([]*domain.Employee, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}
	return m.repo.GetByIDs(ids)
}

package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// EmployeeManager реализует бизнес-логику для сотрудников
type EmployeeManager struct {
	*BaseManager[domain.Employee, *domain.Employee]
}

func NewEmployeeManager(
	repo Repository[domain.Employee, *domain.Employee],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *EmployeeManager {
	return &EmployeeManager{
		BaseManager: NewBaseManager[domain.Employee](repo, logger, txTimeout),
	}
}

// GetByPhone возвращает сотрудника по номеру телефона
func (m *EmployeeManager) GetByPhone(ctx context.Context, phone string) (*domain.Employee, error) {
	employees, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "phone_number", Operator: "=", Value: phone},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByPhone failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "phone", Value: phone},
		)
		return nil, fmt.Errorf("failed to get employee by phone: %w", err)
	}

	if len(employees) == 0 {
		return nil, nil
	}
	return employees[0], nil
}

// GetByUserID возвращает сотрудника по ID пользователя
func (m *EmployeeManager) GetByUserID(ctx context.Context, userID int) (*domain.Employee, error) {
	employees, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "user_id", Operator: "=", Value: userID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByUserID failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "user_id", Value: userID},
		)
		return nil, fmt.Errorf("failed to get employee by user ID: %w", err)
	}

	if len(employees) == 0 {
		return nil, nil
	}
	return employees[0], nil
}

// ListByExperience возвращает сотрудников с опытом работы не менее указанного
func (m *EmployeeManager) ListByExperience(ctx context.Context, minExperience int) ([]*domain.Employee, error) {
	return m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "work_experience", Operator: ">=", Value: minExperience},
		},
		OrderBy: "work_experience DESC",
	})
}

// ListByBirthdayRange возвращает сотрудников с днями рождения в указанном диапазоне
func (m *EmployeeManager) ListByBirthdayRange(ctx context.Context, from, to time.Time) ([]*domain.Employee, error) {
	return m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "birthday", Operator: ">=", Value: from},
			{Field: "birthday", Operator: "<=", Value: to},
		},
		OrderBy: "birthday",
	})
}

// CheckPhoneUnique проверяет уникальность номера телефона
func (m *EmployeeManager) CheckPhoneUnique(ctx context.Context, phone string, excludeID int) (bool, error) {
	employees, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "phone_number", Operator: "=", Value: phone},
			{Field: "employee_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckPhoneUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "phone", Value: phone},
		)
		return false, fmt.Errorf("failed to check phone uniqueness: %w", err)
	}
	return len(employees) == 0, nil
}

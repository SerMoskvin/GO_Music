package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

type EmployeeManager struct {
	*engine.BaseManager[int, domain.Employee, *domain.Employee]
	db *sql.DB
}

func NewEmployeeManager(
	repo db.Repository[domain.Employee, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *EmployeeManager {
	return &EmployeeManager{
		BaseManager: engine.NewBaseManager[int, domain.Employee, *domain.Employee](repo, logger, txTimeout),
		db:          db,
	}
}

// [RU] GetByPhone возвращает сотрудника по номеру телефона <--->
// [ENG] GetByPhone returns an employee by phone number
func (m *EmployeeManager) GetByPhone(ctx context.Context, phone string) (*domain.Employee, error) {
	employees, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "phone_number", Operator: "=", Value: phone},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("GetByPhone failed",
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

// [RU] GetByUserID возвращает сотрудника по ID пользователя <--->
// [ENG] GetByUserID returns an employee by user ID
func (m *EmployeeManager) GetByUserID(ctx context.Context, userID int) (*domain.Employee, error) {
	employees, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "user_id", Operator: "=", Value: userID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("GetByUser ID failed",
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

// [RU] ListByExperience возвращает сотрудников с опытом работы не менее указанного <--->
// [ENG] ListByExperience returns employees with at least the specified work experience
func (m *EmployeeManager) ListByExperience(ctx context.Context, minExperience int) ([]*domain.Employee, error) {
	employees, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "work_experience", Operator: ">=", Value: minExperience},
		},
		OrderBy: "work_experience DESC",
	})
	if err != nil {
		m.Logger.Error("ListByExperience failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "min_experience", Value: minExperience},
		)
		return nil, fmt.Errorf("failed to list employees by experience: %w", err)
	}
	return employees, nil
}

// [RU] ListByBirthdayRange возвращает сотрудников с днями рождения в указанном диапазоне <--->
// [ENG] ListByBirthdayRange returns employees with birthdays in the specified range
func (m *EmployeeManager) ListByBirthdayRange(ctx context.Context, from, to time.Time) ([]*domain.Employee, error) {
	employees, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "birthday", Operator: ">=", Value: from},
			{Field: "birthday", Operator: "<=", Value: to},
		},
		OrderBy: "birthday",
	})
	if err != nil {
		m.Logger.Error("ListByBirthdayRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "from", Value: from},
			logger.Field{Key: "to", Value: to},
		)
		return nil, fmt.Errorf("failed to list employees by birthday range: %w", err)
	}
	return employees, nil
}

// [RU] CheckPhoneUnique проверяет уникальность номера телефона (исключая указанного сотрудника) <--->
// [ENG] CheckPhoneUnique checks the uniqueness of the phone number (excluding the specified employee)
func (m *EmployeeManager) CheckPhoneUnique(ctx context.Context, phone string, excludeID int) (bool, error) {
	employees, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "phone_number", Operator: "=", Value: phone},
			{Field: "employee_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckPhoneUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "phone", Value: phone},
			logger.Field{Key: "exclude_id", Value: excludeID},
		)
		return false, fmt.Errorf("failed to check phone uniqueness: %w", err)
	}
	return len(employees) == 0, nil
}

// [RU] Create создает нового сотрудника с проверкой уникальности телефона <--->
// [ENG] Create creates a new employee with phone uniqueness check
func (m *EmployeeManager) Create(ctx context.Context, employee *domain.Employee) error {
	if err := employee.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee", Value: employee},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckPhoneUnique(ctx, employee.PhoneNumber, 0)
	if err != nil {
		return fmt.Errorf("phone uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("phone number %s already exists", employee.PhoneNumber)
	}

	return m.BaseManager.Create(ctx, employee)
}

// [RU] Update обновляет данные сотрудника с проверкой уникальности телефона <--->
// [ENG] Update updates employee data with phone uniqueness check
func (m *EmployeeManager) Update(ctx context.Context, employee *domain.Employee) error {
	if err := employee.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "employee", Value: employee},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckPhoneUnique(ctx, employee.PhoneNumber, employee.EmployeeID)
	if err != nil {
		return fmt.Errorf("phone uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("phone number %s already exists", employee.PhoneNumber)
	}

	return m.BaseManager.Update(ctx, employee)
}

// [RU] BulkCreate массово создает сотрудников в транзакции <--->
// [ENG] BulkCreate creates multiple employees in a transaction
func (m *EmployeeManager) BulkCreate(ctx context.Context, employees []*domain.Employee) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.Repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, emp := range employees {
		if err := emp.Validate(); err != nil {
			return fmt.Errorf("validation failed for employee %v: %w", emp, err)
		}

		isUnique, err := m.CheckPhoneUnique(ctx, emp.PhoneNumber, 0)
		if err != nil {
			return fmt.Errorf("phone uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("phone number %s already exists", emp.PhoneNumber)
		}

		if err := txRepo.Create(ctx, emp); err != nil {
			return fmt.Errorf("create failed for employee %v: %w", emp, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

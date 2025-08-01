package engine

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// StudentManager реализует бизнес-логику для работы со студентами
type StudentManager struct {
	*BaseManager[int, *domain.Student]
	db *sql.DB
}

func NewStudentManager(
	repo db.Repository[*domain.Student, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentManager {
	return &StudentManager{
		BaseManager: NewBaseManager[int, *domain.Student](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByGroup возвращает студентов указанной группы
func (m *StudentManager) GetByGroup(ctx context.Context, groupID int) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "group_id", Operator: "=", Value: groupID},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.logger.Error("GetByGroup failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group_id", Value: groupID},
		)
		return nil, fmt.Errorf("failed to get students by group: %w", err)
	}
	return DereferenceSlice(students), nil
}

// GetByProgram возвращает студентов по программе обучения
func (m *StudentManager) GetByProgram(ctx context.Context, programID int) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programID},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.logger.Error("GetByProgram failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "program_id", Value: programID},
		)
		return nil, fmt.Errorf("failed to get students by program: %w", err)
	}
	return DereferenceSlice(students), nil
}

// SearchByName ищет студентов по ФИО
func (m *StudentManager) SearchByName(ctx context.Context, query string) ([]*domain.Student, error) {
	repo, ok := m.repo.(interface {
		SearchByName(ctx context.Context, query string) ([]*domain.Student, error)
	})
	if !ok {
		return nil, fmt.Errorf("repository doesn't support SearchByName")
	}

	students, err := repo.SearchByName(ctx, query)
	if err != nil {
		m.logger.Error("SearchByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "query", Value: query},
		)
		return nil, fmt.Errorf("failed to search students by name: %w", err)
	}
	return students, nil
}

// GetByBirthdayRange возвращает студентов в диапазоне дат рождения
func (m *StudentManager) GetByBirthdayRange(ctx context.Context, from, to time.Time) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "birthday", Operator: ">=", Value: from},
			{Field: "birthday", Operator: "<=", Value: to},
		},
		OrderBy: "birthday",
	})
	if err != nil {
		m.logger.Error("GetByBirthdayRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "from", Value: from},
			logger.Field{Key: "to", Value: to},
		)
		return nil, fmt.Errorf("failed to get students by birthday range: %w", err)
	}
	return DereferenceSlice(students), nil
}

// TransferToGroup переводит студента в другую группу
func (m *StudentManager) TransferToGroup(ctx context.Context, studentID, newGroupID int) error {
	studentPtr, err := m.GetByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("failed to get student: %w", err)
	}
	if studentPtr == nil {
		return fmt.Errorf("student not found")
	}

	student := *studentPtr
	student.GroupID = newGroupID

	updatedStudent := &student
	return m.Update(ctx, updatedStudent)
}

// ChangeProgram изменяет программу обучения студента
func (m *StudentManager) ChangeProgram(ctx context.Context, studentID, newProgramID int) error {
	studentPtr, err := m.GetByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("failed to get student: %w", err)
	}
	if studentPtr == nil {
		return fmt.Errorf("student not found")
	}

	student := *studentPtr
	student.MusprogrammID = newProgramID

	updatedStudent := &student
	return m.Update(ctx, updatedStudent)
}

// GetWithUserAccount возвращает студентов с привязанными учетными записями
func (m *StudentManager) GetWithUserAccount(ctx context.Context) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "user_id", Operator: "IS NOT NULL"},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.logger.Error("GetWithUserAccount failed",
			logger.Field{Key: "error", Value: err},
		)
		return nil, fmt.Errorf("failed to get students with user accounts: %w", err)
	}
	return DereferenceSlice(students), nil
}

// CheckPhoneNumberUnique проверяет уникальность номера телефона
func (m *StudentManager) CheckPhoneNumberUnique(ctx context.Context, phone string, excludeID int) (bool, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "phone_number", Operator: "=", Value: phone},
			{Field: "student_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckPhoneNumberUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "phone", Value: phone},
		)
		return false, fmt.Errorf("failed to check phone uniqueness: %w", err)
	}
	return len(students) == 0, nil
}

// Create создает нового студента
func (m *StudentManager) Create(ctx context.Context, student *domain.Student) error {
	if err := student.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student", Value: student},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	if student.PhoneNumber != nil {
		isUnique, err := m.CheckPhoneNumberUnique(ctx, *student.PhoneNumber, 0)
		if err != nil {
			return fmt.Errorf("phone uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("phone number %s already exists", *student.PhoneNumber)
		}
	}

	ptrToStudent := &student
	if err := m.repo.Create(ctx, ptrToStudent); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student", Value: student},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

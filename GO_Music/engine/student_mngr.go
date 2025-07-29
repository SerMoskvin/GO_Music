package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// StudentManager реализует бизнес-логику для работы со студентами
type StudentManager struct {
	*BaseManager[domain.Student, *domain.Student]
}

func NewStudentManager(
	repo Repository[domain.Student, *domain.Student],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentManager {
	return &StudentManager{
		BaseManager: NewBaseManager[domain.Student](repo, logger, txTimeout),
	}
}

// GetByGroup возвращает студентов указанной группы
func (m *StudentManager) GetByGroup(ctx context.Context, groupID int) ([]*domain.Student, error) {
	students, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return students, nil
}

// GetByProgram возвращает студентов по программе обучения
func (m *StudentManager) GetByProgram(ctx context.Context, programID int) ([]*domain.Student, error) {
	students, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return students, nil
}

// SearchByName ищет студентов по ФИО
func (m *StudentManager) SearchByName(ctx context.Context, query string) ([]*domain.Student, error) {
	students, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{
				Field:    "CONCAT(surname, ' ', name, ' ', COALESCE(father_name, ''))",
				Operator: "ILIKE",
				Value:    "%" + query + "%",
			},
		},
		OrderBy: "surname, name",
	})
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
	students, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return students, nil
}

// TransferToGroup переводит студента в другую группу
func (m *StudentManager) TransferToGroup(ctx context.Context, studentID, newGroupID int) error {
	student, err := m.GetByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("failed to get student: %w", err)
	}
	if student == nil {
		return fmt.Errorf("student not found")
	}

	student.GroupID = newGroupID
	return m.Update(ctx, student)
}

// ChangeProgram изменяет программу обучения студента
func (m *StudentManager) ChangeProgram(ctx context.Context, studentID, newProgramID int) error {
	student, err := m.GetByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("failed to get student: %w", err)
	}
	if student == nil {
		return fmt.Errorf("student not found")
	}

	student.MusprogrammID = newProgramID
	return m.Update(ctx, student)
}

// GetWithUserAccount возвращает студентов с привязанными учетными записями
func (m *StudentManager) GetWithUserAccount(ctx context.Context) ([]*domain.Student, error) {
	students, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return students, nil
}

// CheckPhoneNumberUnique проверяет уникальность номера телефона
func (m *StudentManager) CheckPhoneNumberUnique(ctx context.Context, phone string, excludeID int) (bool, error) {
	students, err := m.List(ctx, Filter{
		Conditions: []Condition{
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

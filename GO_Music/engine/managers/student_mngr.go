package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

type StudentManager struct {
	*engine.BaseManager[int, domain.Student, *domain.Student]
	repo *repositories.StudentRepository
	db   *sql.DB
}

func NewStudentManager(
	repo *repositories.StudentRepository,
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *StudentManager {
	return &StudentManager{
		BaseManager: engine.NewBaseManager[int, domain.Student, *domain.Student](repo, logger, txTimeout),
		repo:        repo,
		db:          db,
	}
}

// [RU] GetByGroup возвращает студентов указанной группы <--->
// [ENG] GetByGroup returns students of the specified group
func (m *StudentManager) GetByGroup(ctx context.Context, groupID int) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "group_id", Operator: "=", Value: groupID},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.Logger.Error("GetByGroup failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "group_id", Value: groupID},
		)
		return nil, fmt.Errorf("failed to get students by group: %w", err)
	}
	return students, nil
}

// [RU] GetByProgram возвращает студентов по программе обучения <--->
// [ENG] GetByProgram returns students by study program
func (m *StudentManager) GetByProgram(ctx context.Context, programID int) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "musprogramm_id", Operator: "=", Value: programID},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.Logger.Error("GetByProgram failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "program_id", Value: programID},
		)
		return nil, fmt.Errorf("failed to get students by program: %w", err)
	}
	return students, nil
}

// [RU] SearchByName ищет студентов по ФИО <--->
// [ENG] SearchByName searches for students by full name
func (m *StudentManager) SearchByName(ctx context.Context, query string) ([]*domain.Student, error) {
	students, err := m.repo.SearchByName(ctx, query)
	if err != nil {
		m.Logger.Error("SearchByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "query", Value: query},
		)
		return nil, fmt.Errorf("failed to search students by name: %w", err)
	}
	return students, nil
}

// [RU] GetByBirthdayRange возвращает студентов в диапазоне дат рождения <--->
// [ENG] GetByBirthdayRange returns students in the date of birth range
func (m *StudentManager) GetByBirthdayRange(ctx context.Context, from, to time.Time) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "birthday", Operator: ">=", Value: from},
			{Field: "birthday", Operator: "<=", Value: to},
		},
		OrderBy: "birthday",
	})
	if err != nil {
		m.Logger.Error("GetByBirthdayRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "from", Value: from},
			logger.Field{Key: "to", Value: to},
		)
		return nil, fmt.Errorf("failed to get students by birthday range: %w", err)
	}
	return students, nil
}

// [RU] TransferToGroup переводит студента в другую группу <--->
// [ENG] TransferToGroup transfers a student to another group
func (m *StudentManager) TransferToGroup(ctx context.Context, studentID, newGroupID int) error {
	studentPtr, err := m.GetByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("failed to get student: %w", err)
	}
	if studentPtr == nil {
		return fmt.Errorf("student not found")
	}

	studentPtr.GroupID = newGroupID
	return m.Update(ctx, studentPtr)
}

// [RU] ChangeProgram изменяет программу обучения студента <--->
// [ENG] ChangeProgram changes the student's study program
func (m *StudentManager) ChangeProgram(ctx context.Context, studentID, newProgramID int) error {
	studentPtr, err := m.GetByID(ctx, studentID)
	if err != nil {
		return fmt.Errorf("failed to get student: %w", err)
	}
	if studentPtr == nil {
		return fmt.Errorf("student not found")
	}

	studentPtr.MusprogrammID = newProgramID
	return m.Update(ctx, studentPtr)
}

// [RU] GetWithUserAccount возвращает студентов с привязанными учетными записями <--->
// [ENG] GetWithUserAccount returns students with associated user accounts
func (m *StudentManager) GetWithUserAccount(ctx context.Context) ([]*domain.Student, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "user_id", Operator: "IS NOT NULL"},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.Logger.Error("GetWithUser Account failed",
			logger.Field{Key: "error", Value: err},
		)
		return nil, fmt.Errorf("failed to get students with user accounts: %w", err)
	}
	return students, nil
}

// [RU] CheckPhoneNumberUnique проверяет уникальность номера телефона <--->
// [ENG] CheckPhoneNumberUnique checks the uniqueness of the phone number
func (m *StudentManager) CheckPhoneNumberUnique(ctx context.Context, phone string, excludeID int) (bool, error) {
	students, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "phone_number", Operator: "=", Value: phone},
			{Field: "student_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckPhoneNumberUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "phone", Value: phone},
		)
		return false, fmt.Errorf("failed to check phone uniqueness: %w", err)
	}
	return len(students) == 0, nil
}

// [RU] Create создает нового студента <--->
// [ENG] Create creates a new student
func (m *StudentManager) Create(ctx context.Context, student *domain.Student) error {
	if err := student.Validate(); err != nil {
		m.Logger.Error("Validation failed",
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

	if err := m.Repo.Create(ctx, student); err != nil {
		m.Logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student", Value: student},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// [RU] Update обновляет данные студента <--->
// [ENG] Update updates the student's data
func (m *StudentManager) Update(ctx context.Context, student *domain.Student) error {
	if err := student.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student", Value: student},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	if student.PhoneNumber != nil {
		isUnique, err := m.CheckPhoneNumberUnique(ctx, *student.PhoneNumber, student.StudentID)
		if err != nil {
			return fmt.Errorf("phone uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("phone number %s already exists", *student.PhoneNumber)
		}
	}

	if err := m.Repo.Update(ctx, student); err != nil {
		m.Logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "student", Value: student},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// [RU] BulkCreate массово создает студентов в транзакции <--->
// [ENG] BulkCreate creates multiple students in a transaction
func (m *StudentManager) BulkCreate(ctx context.Context, students []*domain.Student) error {
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

	for _, student := range students {
		if err := student.Validate(); err != nil {
			return fmt.Errorf("validation failed for student %v: %w", student, err)
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

		if err := txRepo.Create(ctx, student); err != nil {
			return fmt.Errorf("create failed for student %v: %w", student, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

package engine

import (
	"GO_Music/domain"
	"errors"

	"github.com/SerMoskvin/access"
)



type StudentRepository interface {
	Repository[domain.Student, *domain.Student]
	GetByTeacherID(teacherID int) ([]*domain.Student, error)
}

type StudentManager struct {
	*BaseManager[domain.Student, *domain.Student]
	groupRepo   StudyGroupRepository
	permissions map[string]access.RolePermissions
}

func NewStudentManager(
	repo StudentRepository,
	groupRepo StudyGroupRepository,
	permissions map[string]access.RolePermissions,
) *StudentManager {
	return &StudentManager{
		BaseManager: NewBaseManager[domain.Student, *domain.Student](repo),
		groupRepo:   groupRepo,
		permissions: permissions,
	}
}

func (m *StudentManager) GetByIDWithAccess(id int, userID int, userRole string) (*domain.Student, error) {
	rolePerms, ok := m.permissions[userRole] // Используем поле
	if !ok {
		return nil, errors.New("неизвестная роль")
	}

	if rolePerms.OwnRecordsOnly && id != userID {
		return nil, errors.New("доступ запрещён")
	}

	return m.GetByID(id)
}

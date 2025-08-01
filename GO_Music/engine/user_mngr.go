package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/domain"

	"github.com/SerMoskvin/access"
	"github.com/SerMoskvin/logger"
	"github.com/dgrijalva/jwt-go"
)

// UserManager реализует бизнес-логику для работы с пользователями
type UserManager struct {
	*BaseManager[int, *domain.User]
	db   *sql.DB
	auth *access.Authenticator
}

func NewUserManager(
	repo db.Repository[*domain.User, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
	auth *access.Authenticator,
) *UserManager {
	return &UserManager{
		BaseManager: NewBaseManager[int, *domain.User](repo, logger, txTimeout),
		db:          db,
		auth:        auth,
	}
}

// Register создает нового пользователя с хешированным паролем
func (m *UserManager) Register(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Проверка уникальности логина
	isUnique, err := m.CheckLoginUnique(ctx, user.Login, 0)
	if err != nil {
		return fmt.Errorf("login uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("login %s already exists", user.Login)
	}

	// Хеширование пароля
	hashedPassword, err := m.auth.PasswordHasher.HashPassword(user.Password)
	if err != nil {
		m.logger.Error("Failed to hash password",
			logger.Field{Key: "error", Value: err},
		)
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Установка даты регистрации
	user.RegistrationDate = time.Now()

	ptrToUser := &user
	return m.repo.Create(ctx, ptrToUser)
}

// Login выполняет аутентификацию пользователя и возвращает JWT-токен
func (m *UserManager) Login(ctx context.Context, login, password string) (string, error) {
	users, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "login", Operator: "=", Value: login},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("Login failed - user search error",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "login", Value: login},
		)
		return "", fmt.Errorf("authentication failed")
	}

	if len(users) == 0 {
		m.logger.Warn("Login failed - user not found",
			logger.Field{Key: "login", Value: login},
		)
		return "", fmt.Errorf("authentication failed")
	}

	user := *users[0]

	// Проверка пароля
	if !m.auth.PasswordHasher.CheckPasswordHash(password, user.Password) {
		m.logger.Warn("Login failed - invalid password",
			logger.Field{Key: "login", Value: login},
		)
		return "", fmt.Errorf("authentication failed")
	}

	// Генерация токена
	token, err := m.auth.JwtService.GenerateJWT(user.UserID, user.Login, user.Role)
	if err != nil {
		m.logger.Error("Login failed - token generation error",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "user_id", Value: user.UserID},
		)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// GetCurrentUser возвращает данные текущего аутентифицированного пользователя
func (m *UserManager) GetCurrentUser(ctx context.Context) (*domain.User, error) {
	claims, ok := ctx.Value(access.UserClaimsKey).(jwt.MapClaims)
	if !ok {
		return nil, errors.New("authentication required")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userPtr, err := m.GetByID(ctx, int(userID))
	if err != nil {
		return nil, err
	}
	if userPtr == nil {
		return nil, errors.New("user not found")
	}
	return *userPtr, nil
}

// ChangePassword изменяет пароль пользователя
func (m *UserManager) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	userPtr, err := m.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userPtr == nil {
		return fmt.Errorf("user not found")
	}

	user := *userPtr

	if !m.auth.PasswordHasher.CheckPasswordHash(oldPassword, user.Password) {
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := m.auth.PasswordHasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.Password = hashedPassword
	ptrToUser := &user
	return m.Update(ctx, ptrToUser)
}

// GetByRole возвращает пользователей по роли
func (m *UserManager) GetByRole(ctx context.Context, role string) ([]*domain.User, error) {
	users, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "role", Operator: "=", Value: role},
		},
		OrderBy: "surname, name",
	})
	if err != nil {
		m.logger.Error("GetByRole failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "role", Value: role},
		)
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	return DereferenceSlice(users), nil
}

// SearchByNames ищет пользователей по ФИО
func (m *UserManager) SearchByNames(ctx context.Context, query string) ([]*domain.User, error) {
	repo, ok := m.repo.(interface {
		SearchByName(ctx context.Context, query string) ([]*domain.User, error)
	})
	if !ok {
		return nil, fmt.Errorf("repository doesn't support SearchByName")
	}

	users, err := repo.SearchByName(ctx, query)
	if err != nil {
		m.logger.Error("SearchByNames failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "query", Value: query},
		)
		return nil, fmt.Errorf("failed to search users by names: %w", err)
	}
	return users, nil
}

// CheckLoginUnique проверяет уникальность логина
func (m *UserManager) CheckLoginUnique(ctx context.Context, login string, excludeID int) (bool, error) {
	users, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "login", Operator: "=", Value: login},
			{Field: "user_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckLoginUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "login", Value: login},
		)
		return false, fmt.Errorf("failed to check login uniqueness: %w", err)
	}
	return len(users) == 0, nil
}

// UpdateProfile обновляет профиль пользователя
func (m *UserManager) UpdateProfile(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Проверка уникальности логина
	isUnique, err := m.CheckLoginUnique(ctx, user.Login, user.UserID)
	if err != nil {
		return fmt.Errorf("login uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("login %s already exists", user.Login)
	}

	ptrToUser := &user
	return m.Update(ctx, ptrToUser)
}

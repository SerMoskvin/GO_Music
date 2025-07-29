package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/access"

	"github.com/SerMoskvin/logger"
	"github.com/dgrijalva/jwt-go"
)

// UserManager реализует бизнес-логику для работы с пользователями
type UserManager struct {
	*BaseManager[domain.User, *domain.User]
	auth *access.Authenticator
}

func NewUserManager(
	repo Repository[domain.User, *domain.User],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
	auth *access.Authenticator,
) *UserManager {
	return &UserManager{
		BaseManager: NewBaseManager[domain.User](repo, logger, txTimeout),
		auth:        auth,
	}
}

// Register создает нового пользователя с хешированным паролем
func (m *UserManager) Register(ctx context.Context, user *domain.User) error {
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

	return m.Create(ctx, user)
}

// Login выполняет аутентификацию пользователя и возвращает JWT-токен
func (m *UserManager) Login(ctx context.Context, login, password string) (string, error) {
	// Находим пользователя по логину
	users, err := m.List(ctx, Filter{
		Conditions: []Condition{
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

	user := users[0]

	// Проверяем пароль
	if !m.auth.PasswordHasher.CheckPasswordHash(password, user.Password) {
		m.logger.Warn("Login failed - invalid password",
			logger.Field{Key: "login", Value: login},
		)
		return "", fmt.Errorf("authentication failed")
	}

	// Генерируем JWT-токен
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

	return m.GetByID(ctx, int(userID))
}

// ChangePassword изменяет пароль пользователя с проверкой старого
func (m *UserManager) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := m.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !m.auth.PasswordHasher.CheckPasswordHash(oldPassword, user.Password) {
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := m.auth.PasswordHasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.Password = hashedPassword
	return m.Update(ctx, user)
}

// GetByRole возвращает пользователей по роли
func (m *UserManager) GetByRole(ctx context.Context, role string) ([]*domain.User, error) {
	users, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return users, nil
}

// SearchByNames ищет пользователей по ФИО
func (m *UserManager) SearchByNames(ctx context.Context, query string) ([]*domain.User, error) {
	users, err := m.List(ctx, Filter{
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
	users, err := m.List(ctx, Filter{
		Conditions: []Condition{
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

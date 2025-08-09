package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	"github.com/SerMoskvin/access"
	"github.com/SerMoskvin/logger"
	"github.com/dgrijalva/jwt-go"
)

type UserManager struct {
	*BaseManager[int, domain.User, *domain.User]
	repo *repositories.UserRepository
	db   *sql.DB
	auth *access.Authenticator
}

func NewUserManager(
	repo *repositories.UserRepository,
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
	auth *access.Authenticator,
) *UserManager {
	return &UserManager{
		BaseManager: NewBaseManager[int, domain.User, *domain.User](repo, logger, txTimeout),
		repo:        repo,
		db:          db,
		auth:        auth,
	}
}

// [RU] Register создает нового пользователя с хешированным паролем <--->
// [ENG] Register creates a new user with a hashed password
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

	return m.repo.Create(ctx, user)
}

// [RU] Login выполняет аутентификацию пользователя и возвращает JWT-токен <--->
// [ENG] Login authenticates the user and returns a JWT token
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

// [RU] GetCurrentUser  возвращает данные текущего аутентифицированного пользователя <--->
// [ENG] GetCurrentUser  returns the data of the currently authenticated user
func (m *UserManager) GetCurrentUser(ctx context.Context) (*domain.User, error) {
	claims, ok := ctx.Value(access.UserClaimsKey).(jwt.MapClaims)
	if !ok {
		return nil, errors.New("authentication required")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	userID := int(userIDFloat)

	user, err := m.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// [RU] ChangePassword изменяет пароль пользователя <--->
// [ENG] ChangePassword changes the user's password
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

// [RU] GetByRole возвращает пользователей по роли <--->
// [ENG] GetByRole returns users by role
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
	return users, nil
}

// [RU] SearchByNames ищет пользователей по ФИО <--->
// [ENG] SearchByNames searches for users by full name
func (m *UserManager) SearchByNames(ctx context.Context, query string) ([]*domain.User, error) {
	users, err := m.repo.SearchByName(ctx, query)
	if err != nil {
		m.logger.Error("SearchByNames failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "query", Value: query},
		)
		return nil, fmt.Errorf("failed to search users by names: %w", err)
	}

	for _, user := range users {
		if user.Image == nil {
			user.Image = DefaultImage // Подставляем изображение по умолчанию
		}
	}

	return users, nil
}

// [RU] CheckLoginUnique проверяет уникальность логина <--->
// [ENG] CheckLoginUnique checks the uniqueness of the login
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

// [RU] UpdateProfile обновляет профиль пользователя <--->
// [ENG] UpdateProfile updates the user's profile
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

	return m.Update(ctx, user)
}

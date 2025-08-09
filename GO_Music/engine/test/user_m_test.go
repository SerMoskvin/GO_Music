package engine_test

import (
	"context"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/access"
	"github.com/SerMoskvin/logger"
	"github.com/stretchr/testify/assert"
)

func TestUserManager_AllMethods(t *testing.T) {
	// Пути к конфигам
	cfgPathDB := "../../config/DB_config.yml"
	cfgPathLog := "../../config/logger_config.yml"
	cfgPathAccess := "../../config/access_config.yml"

	// Загружаем DB конфиг и инициализируем базу
	cfgDB, err := config.LoadDBConfig(cfgPathDB)
	if err != nil {
		t.Fatalf("failed to load db config: %v", err)
	}

	sqlDB, err := db.InitPostgresDB(cfgDB)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаём репозиторий
	repo := repositories.NewUserRepository(sqlDB)

	// Создаём логгер
	levelLogger, err := logger.NewLevel(cfgPathLog)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	// Создаём аутентификатор через ваш NewAuthenticator с конфигом
	authenticator, err := access.NewAuthenticator(cfgPathAccess)
	if err != nil {
		t.Fatalf("failed to create authenticator: %v", err)
	}

	// Создаём UserManager
	mgr := engine.NewUserManager(repo, sqlDB, levelLogger, 5*time.Second, authenticator)

	// Тестовые данные
	testUser := &domain.User{
		Login:            "testuser123",
		Password:         "password123",
		Name:             "Иван",
		Surname:          "Иванов",
		Role:             "student",
		RegistrationDate: domain.ParseDMY("15.05.2025"),
		Email:            "Test_test@mail.ru",
	}

	// Регистрация
	t.Run("Register", func(t *testing.T) {
		err := mgr.Register(ctx, testUser)
		if err != nil {
			levelLogger.Error("Register failed", logger.String("error", err.Error()))
			t.Fatalf("Register failed: %v", err)
		}
		assert.NotZero(t, testUser.UserID)
	})

	// Логин
	t.Run("Login", func(t *testing.T) {
		token, err := mgr.Login(ctx, testUser.Login, "password123")
		if err != nil {
			levelLogger.Error("Login failed", logger.String("error", err.Error()))
			t.Fatalf("Login failed: %v", err)
		}
		assert.NotEmpty(t, token)
	})

	// Получение пользователей по роли
	t.Run("GetByRole", func(t *testing.T) {
		users, err := mgr.GetByRole(ctx, testUser.Role)
		if err != nil {
			levelLogger.Error("GetByRole failed", logger.String("error", err.Error()), logger.String("role", testUser.Role))
			t.Fatalf("GetByRole failed: %v", err)
		}
		assert.NotEmpty(t, users)
	})

	// Поиск по ФИО
	t.Run("SearchByNames", func(t *testing.T) {
		users, err := mgr.SearchByNames(ctx, "Иван")
		if err != nil {
			levelLogger.Error("SearchByNames failed", logger.String("error", err.Error()), logger.String("query", "Иван"))
			t.Fatalf("SearchByNames failed: %v", err)
		}
		assert.NotEmpty(t, users)
	})

	// Проверка уникальности логина
	t.Run("CheckLoginUnique", func(t *testing.T) {
		unique, err := mgr.CheckLoginUnique(ctx, testUser.Login, testUser.UserID)
		if err != nil {
			levelLogger.Error("CheckLoginUnique failed", logger.String("error", err.Error()), logger.String("login", testUser.Login))
			t.Fatalf("CheckLoginUnique failed: %v", err)
		}
		assert.True(t, unique)
	})

	// Смена пароля
	t.Run("ChangePassword", func(t *testing.T) {
		newPassword := "newpassword123"
		err := mgr.ChangePassword(ctx, testUser.UserID, "password123", newPassword)
		if err != nil {
			levelLogger.Error("ChangePassword failed", logger.String("error", err.Error()))
			t.Fatalf("ChangePassword failed: %v", err)
		}

		// Старый пароль не должен работать
		_, err = mgr.Login(ctx, testUser.Login, "password123")
		assert.Error(t, err)

		// Новый пароль должен работать
		token, err := mgr.Login(ctx, testUser.Login, newPassword)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	// Обновление профиля
	t.Run("UpdateProfile", func(t *testing.T) {
		testUser.Name = "Пётр"
		testUser.Surname = "Петров"
		err := mgr.UpdateProfile(ctx, testUser)
		if err != nil {
			levelLogger.Error("UpdateProfile failed", logger.String("error", err.Error()))
			t.Fatalf("UpdateProfile failed: %v", err)
		}

		updatedUser, err := mgr.GetByID(ctx, testUser.UserID)
		if err != nil {
			levelLogger.Error("GetByID after UpdateProfile failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID after UpdateProfile failed: %v", err)
		}
		assert.Equal(t, "Пётр", updatedUser.Name)
		assert.Equal(t, "Петров", updatedUser.Surname)
	})

	// Удаление пользователя
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, testUser.UserID)
		if err != nil {
			levelLogger.Error("Delete failed", logger.String("error", err.Error()))
			t.Fatalf("Delete failed: %v", err)
		}

		exists, err := repo.Exists(ctx, testUser.UserID)
		if err != nil {
			levelLogger.Error("Exists after delete failed", logger.String("error", err.Error()))
			t.Fatalf("Exists after delete failed: %v", err)
		}
		assert.False(t, exists)
	})

	if !t.Failed() {
		levelLogger.Info("All UserManager tests passed successfully")
		t.Log("All UserManager tests passed successfully")
	}
}

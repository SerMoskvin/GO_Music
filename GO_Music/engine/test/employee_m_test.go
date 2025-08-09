package engine_test

import (
	"context"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"
	"GO_Music/engine/managers"

	"github.com/SerMoskvin/logger"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeManager_AllMethods(t *testing.T) {
	// Загрузка конфигурации
	cfgPath_DB := "../../config/DB_config.yml"
	cfgPath_Log := "../../config/logger_config.yml"
	cfg, err := config.LoadDBConfig(cfgPath_DB)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	sqlDB, err := db.InitPostgresDB(cfg)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo := repositories.NewEmployeeRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := managers.NewEmployeeManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	UserID := 5
	testEmployee := &domain.Employee{
		PhoneNumber:    "79109532210",
		UserID:         &UserID,
		WorkExperience: 5,
		Birthday:       time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		Job:            "Testovik",
		Surname:        "Тестов",
		Name:           "Аркадий",
		FatherName:     nil,
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testEmployee)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		if testEmployee.EmployeeID == 0 {
			levelLogger.Error("Expected ID to be set after Create")
			t.Error("Expected ID to be set after Create")
		}
		levelLogger.Info("Created employee", logger.Int("ID", testEmployee.EmployeeID))
	})

	t.Run("GetByPhone", func(t *testing.T) {
		emp, err := mgr.GetByPhone(ctx, testEmployee.PhoneNumber)
		assert.NoError(t, err)
		assert.NotNil(t, emp)
		if emp != nil {
			assert.Equal(t, testEmployee.PhoneNumber, emp.PhoneNumber)
		}
	})

	t.Run("GetByUserID", func(t *testing.T) {
		emp, err := mgr.GetByUserID(ctx, *testEmployee.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, emp)
		if emp != nil {
			assert.Equal(t, testEmployee.UserID, emp.UserID)
		}
	})

	t.Run("CheckPhoneUnique", func(t *testing.T) {
		isUnique, err := mgr.CheckPhoneUnique(ctx, testEmployee.PhoneNumber, testEmployee.EmployeeID)
		assert.NoError(t, err)
		assert.True(t, isUnique, "Phone should be unique excluding current employee")

		isUnique, err = mgr.CheckPhoneUnique(ctx, testEmployee.PhoneNumber, 0)
		assert.NoError(t, err)
		assert.False(t, isUnique, "Phone should not be unique without exclusion")

		isUnique, err = mgr.CheckPhoneUnique(ctx, "nonexistentphone", 0)
		assert.NoError(t, err)
		assert.True(t, isUnique, "Nonexistent phone should be unique")
	})

	t.Run("ListByExperience", func(t *testing.T) {
		emps, err := mgr.ListByExperience(ctx, 3)
		assert.NoError(t, err)
		assert.NotEmpty(t, emps)

		found := false
		for _, e := range emps {
			if e.EmployeeID == testEmployee.EmployeeID {
				found = true
				assert.GreaterOrEqual(t, e.WorkExperience, 3)
			}
		}
		assert.True(t, found, "Created employee should be in ListByExperience result")
	})

	t.Run("ListByBirthdayRange", func(t *testing.T) {
		from := time.Date(1989, 12, 31, 0, 0, 0, 0, time.UTC)
		to := time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC)
		emps, err := mgr.ListByBirthdayRange(ctx, from, to)
		assert.NoError(t, err)
		assert.NotEmpty(t, emps)

		found := false
		for _, e := range emps {
			if e.EmployeeID == testEmployee.EmployeeID {
				found = true
				assert.True(t, !e.Birthday.Before(from) && !e.Birthday.After(to))
			}
		}
		assert.True(t, found, "Created employee should be in ListByBirthdayRange result")
	})

	t.Run("Update", func(t *testing.T) {
		newPhone := "79109345566"
		testEmployee.PhoneNumber = newPhone
		err := mgr.Update(ctx, testEmployee)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		emp, err := mgr.GetByPhone(ctx, newPhone)
		assert.NoError(t, err)
		assert.NotNil(t, emp)
		if emp != nil {
			assert.Equal(t, newPhone, emp.PhoneNumber)
		}
	})

	t.Run("Create Duplicate Phone", func(t *testing.T) {
		dupEmployee := &domain.Employee{
			PhoneNumber:    testEmployee.PhoneNumber,
			UserID:         &UserID,
			WorkExperience: 2,
			Birthday:       time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC),
			Job:            "Бухгалтер",
			Surname:        "Тестов",
			Name:           "Аркадий",
			FatherName:     nil,
		}
		err := mgr.Create(ctx, dupEmployee)
		assert.Error(t, err)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		UserID_1, UserID_2 := 6, 7
		FatherName := "Владимирович"
		emps := []*domain.Employee{
			{
				PhoneNumber:    "79109887645",
				UserID:         &UserID_1,
				WorkExperience: 1,
				Birthday:       time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC),
				Job:            "Тестовик_1",
				Surname:        "Loxov",
				Name:           "Vladimir",
				FatherName:     &FatherName,
			},
			{
				PhoneNumber:    "79109887645",
				UserID:         &UserID_2,
				WorkExperience: 4,
				Birthday:       time.Date(1988, 8, 8, 0, 0, 0, 0, time.UTC),
				Job:            "Тестовик_2",
				Surname:        "Иванонв",
				Name:           "Иван",
				FatherName:     nil,
			},
		}

		err := mgr.BulkCreate(ctx, emps)
		assert.NoError(t, err)

		for _, e := range emps {
			assert.NotZero(t, e.EmployeeID)
		}
	})

	if !t.Failed() {
		levelLogger.Info("All EmployeeManager tests passed successfully")
		t.Log("All EmployeeManager tests passed successfully")
	}
}

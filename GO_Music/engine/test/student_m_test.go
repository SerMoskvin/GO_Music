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

	"github.com/SerMoskvin/logger"
	"github.com/stretchr/testify/assert"
)

func TestStudentManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := engine.NewStudentManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные

	testStudent := &domain.Student{
		UserID:        ptrNumber(3),
		Name:          "Иван",
		Surname:       "Иванов",
		GroupID:       1,
		MusprogrammID: 1,
		Birthday:      domain.ParseDMY("15.12.2006"),
		PhoneNumber:   ptrString("79109546785"),
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testStudent)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		assert.NotZero(t, testStudent.StudentID)
	})

	t.Run("GetByID", func(t *testing.T) {
		studentPtr, err := mgr.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("GetByID failed", logger.String("error", err.Error()), logger.Int("studentID", testStudent.StudentID))
			t.Fatalf("GetByID failed: %v", err)
		}
		if studentPtr == nil {
			levelLogger.Error("Expected student to be found", logger.Int("studentID", testStudent.StudentID))
			t.Fatal("Expected student to be found")
		}
		assert.Equal(t, testStudent.Name, studentPtr.Name)
	})

	t.Run("GetByGroup", func(t *testing.T) {
		students, err := mgr.GetByGroup(ctx, testStudent.GroupID)
		if err != nil {
			levelLogger.Error("GetByGroup failed", logger.String("error", err.Error()), logger.Int("groupID", testStudent.GroupID))
			t.Fatalf("GetByGroup failed: %v", err)
		}
		assert.NotEmpty(t, students)
	})

	t.Run("GetByProgram", func(t *testing.T) {
		students, err := mgr.GetByProgram(ctx, testStudent.MusprogrammID)
		if err != nil {
			levelLogger.Error("GetByProgram failed", logger.String("error", err.Error()), logger.Int("programID", testStudent.MusprogrammID))
			t.Fatalf("GetByProgram failed: %v", err)
		}
		assert.NotEmpty(t, students)
	})

	t.Run("SearchByName", func(t *testing.T) {
		query := "Иван"
		students, err := mgr.SearchByName(ctx, query)
		if err != nil {
			levelLogger.Error("SearchByName failed", logger.String("error", err.Error()), logger.String("query", query))
			t.Fatalf("SearchByName failed: %v", err)
		}
		assert.NotEmpty(t, students)
	})

	t.Run("GetByBirthdayRange", func(t *testing.T) {
		from := domain.ParseDMY("01.01.2005")
		to := domain.ParseDMY("01.01.2007")
		students, err := mgr.GetByBirthdayRange(ctx, from, to)
		if err != nil {
			levelLogger.Error("GetByBirthdayRange failed", logger.String("error", err.Error()), logger.String("from", from.Format(time.RFC3339)), logger.String("to", to.Format(time.RFC3339)))
			t.Fatalf("GetByBirthdayRange failed: %v", err)
		}
		assert.NotEmpty(t, students)
	})

	t.Run("TransferToGroup", func(t *testing.T) {
		newGroupID := 2
		err := mgr.TransferToGroup(ctx, testStudent.StudentID, newGroupID)
		if err != nil {
			levelLogger.Error("TransferToGroup failed", logger.String("error", err.Error()), logger.Int("studentID", testStudent.StudentID), logger.Int("newGroupID", newGroupID))
			t.Fatalf("TransferToGroup failed: %v", err)
		}

		updatedStudent, err := mgr.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("GetByID after TransferToGroup failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID after TransferToGroup failed: %v", err)
		}
		assert.Equal(t, newGroupID, updatedStudent.GroupID)
	})

	t.Run("ChangeProgram", func(t *testing.T) {
		newProgramID := 3
		err := mgr.ChangeProgram(ctx, testStudent.StudentID, newProgramID)
		if err != nil {
			levelLogger.Error("ChangeProgram failed", logger.String("error", err.Error()), logger.Int("studentID", testStudent.StudentID), logger.Int("newProgramID", newProgramID))
			t.Fatalf("ChangeProgram failed: %v", err)
		}

		updatedStudent, err := mgr.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("GetByID after ChangeProgram failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID after ChangeProgram failed: %v", err)
		}
		assert.Equal(t, newProgramID, updatedStudent.MusprogrammID)
	})

	t.Run("GetWithUserAccount", func(t *testing.T) {
		students, err := mgr.GetWithUserAccount(ctx)
		if err != nil {
			levelLogger.Error("GetWithUserAccount failed", logger.String("error", err.Error()))
			t.Fatalf("GetWithUserAccount failed: %v", err)
		}
		assert.NotNil(t, students)
	})

	t.Run("CheckPhoneNumberUnique", func(t *testing.T) {
		unique, err := mgr.CheckPhoneNumberUnique(ctx, *testStudent.PhoneNumber, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("CheckPhoneNumberUnique failed", logger.String("error", err.Error()), logger.String("phone", *testStudent.PhoneNumber))
			t.Fatalf("CheckPhoneNumberUnique failed: %v", err)
		}
		assert.True(t, unique)
	})

	t.Run("Update", func(t *testing.T) {
		testStudent.Name = "Петр"
		err := mgr.Update(ctx, testStudent)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		updatedStudent, err := mgr.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("GetByID after Update failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID after Update failed: %v", err)
		}
		assert.Equal(t, "Петр", updatedStudent.Name)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		students := []*domain.Student{
			{
				UserID:        ptrNumber(6),
				Name:          "Алексей",
				Surname:       "Сидоров",
				GroupID:       1,
				MusprogrammID: 1,
				Birthday:      domain.ParseDMY("10.10.2010"),
				PhoneNumber:   ptrString("79881043122"),
			},
			{
				UserID:        ptrNumber(7),
				Name:          "Мария",
				Surname:       "Петрова",
				GroupID:       2,
				MusprogrammID: 2,
				Birthday:      domain.ParseDMY("12.12.2012"),
				PhoneNumber:   ptrString("79109887676"),
			},
		}

		err := mgr.BulkCreate(ctx, students)
		if err != nil {
			levelLogger.Error("BulkCreate failed", logger.String("error", err.Error()))
			t.Fatalf("BulkCreate failed: %v", err)
		}

		for _, s := range students {
			exists, err := repo.Exists(ctx, s.StudentID)
			if err != nil {
				levelLogger.Error("Exists check failed after BulkCreate", logger.String("error", err.Error()), logger.Int("studentID", s.StudentID))
				t.Fatalf("Exists check failed after BulkCreate: %v", err)
			}
			assert.True(t, exists, "Expected student to exist after BulkCreate")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("Delete failed", logger.String("error", err.Error()))
			t.Fatalf("Delete failed: %v", err)
		}

		exists, err := repo.Exists(ctx, testStudent.StudentID)
		if err != nil {
			levelLogger.Error("Exists after delete failed", logger.String("error", err.Error()))
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			levelLogger.Error("Expected student to be deleted")
			t.Error("Expected student to be deleted")
		}
	})

	if !t.Failed() {
		levelLogger.Info("All StudentManager tests passed successfully")
		t.Log("All StudentManager tests passed successfully")
	}
}

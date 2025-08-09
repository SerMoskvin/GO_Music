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

func TestSubjectManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewSubjectRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := engine.NewSubjectManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testSubject := &domain.Subject{
		SubjectName: "Музыка",
		SubjectType: "Теория",
		ShortDesc:   "Изучение теории музыки",
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testSubject)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		assert.NotZero(t, testSubject.SubjectID)
	})

	t.Run("GetByType", func(t *testing.T) {
		subjects, err := mgr.GetByType(ctx, testSubject.SubjectType)
		if err != nil {
			levelLogger.Error("GetByType failed", logger.String("error", err.Error()), logger.String("type", testSubject.SubjectType))
			t.Fatalf("GetByType failed: %v", err)
		}
		assert.NotEmpty(t, subjects)
	})

	t.Run("SearchByName", func(t *testing.T) {
		subjects, err := mgr.SearchByName(ctx, testSubject.SubjectName)
		if err != nil {
			levelLogger.Error("SearchByName failed", logger.String("error", err.Error()), logger.String("name", testSubject.SubjectName))
			t.Fatalf("SearchByName failed: %v", err)
		}
		assert.NotEmpty(t, subjects)
	})

	t.Run("GetByDescription", func(t *testing.T) {
		subjects, err := mgr.GetByDescription(ctx, "теории")
		if err != nil {
			levelLogger.Error("GetByDescription failed", logger.String("error", err.Error()), logger.String("keyword", "теории"))
			t.Fatalf("GetByDescription failed: %v", err)
		}
		assert.NotEmpty(t, subjects)
	})

	t.Run("CheckNameUnique", func(t *testing.T) {
		unique, err := mgr.CheckNameUnique(ctx, testSubject.SubjectName, testSubject.SubjectID)
		if err != nil {
			levelLogger.Error("CheckNameUnique failed", logger.String("error", err.Error()), logger.String("name", testSubject.SubjectName))
			t.Fatalf("CheckNameUnique failed: %v", err)
		}
		assert.True(t, unique)
	})

	t.Run("Update", func(t *testing.T) {
		testSubject.SubjectName = "Музыка и искусство"
		err := mgr.Update(ctx, testSubject)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		updatedSubject, err := mgr.GetByType(ctx, testSubject.SubjectType)
		if err != nil {
			levelLogger.Error("GetByType after Update failed", logger.String("error", err.Error()))
			t.Fatalf("GetByType after Update failed: %v", err)
		}
		assert.Equal(t, "Музыка и искусство", updatedSubject[0].SubjectName)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		subjects := []*domain.Subject{
			{
				SubjectName: "Test",
				SubjectType: "Теория",
				ShortDesc:   "Изучение истории музыки",
			},
			{
				SubjectName: "Another_test",
				SubjectType: "Практика",
				ShortDesc:   "Изучение практики музыки",
			},
		}

		err := mgr.BulkCreate(ctx, subjects)
		if err != nil {
			levelLogger.Error("BulkCreate failed", logger.String("error", err.Error()))
			t.Fatalf("BulkCreate failed: %v", err)
		}

		for _, s := range subjects {
			exists, err := repo.Exists(ctx, s.SubjectID)
			if err != nil {
				levelLogger.Error("Exists check failed after BulkCreate", logger.String("error", err.Error()), logger.Int("subjectID", s.SubjectID))
				t.Fatalf("Exists check failed after BulkCreate: %v", err)
			}
			assert.True(t, exists, "Expected subject to exist after BulkCreate")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, testSubject.SubjectID)
		if err != nil {
			levelLogger.Error("Delete failed", logger.String("error", err.Error()))
			t.Fatalf("Delete failed: %v", err)
		}

		exists, err := repo.Exists(ctx, testSubject.SubjectID)
		if err != nil {
			levelLogger.Error("Exists after delete failed", logger.String("error", err.Error()))
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			levelLogger.Error("Expected subject to be deleted")
			t.Error("Expected subject to be deleted")
		}
	})

	if !t.Failed() {
		levelLogger.Info("All SubjectManager tests passed successfully")
		t.Log("All SubjectManager tests passed successfully")
	}
}

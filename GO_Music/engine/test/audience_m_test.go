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
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAudienceManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewAudienceRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := managers.NewAudienceManager(repo, levelLogger, 5*time.Second)

	// Тестовые данные
	testAudience := &domain.Audience{
		Name:        "Тестовая запись",
		AudinType:   "Тестовый тип",
		AudinNumber: "A-101",
		Capacity:    22,
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testAudience)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		if testAudience.AudienceID == 0 {
			levelLogger.Error("Expected AudienceID to be set after Create")
			t.Error("Expected AudienceID to be set after Create")
		}
		levelLogger.Info("Created audience", logger.Int("ID", testAudience.AudienceID), logger.String("Number", testAudience.AudinNumber), logger.Int("Capacity", testAudience.Capacity))
	})

	t.Run("GetByNumber", func(t *testing.T) {
		aud, err := mgr.GetByNumber(ctx, testAudience.AudinNumber)
		assert.NoError(t, err)
		assert.NotNil(t, aud)
		if aud != nil {
			assert.Equal(t, testAudience.AudinNumber, aud.AudinNumber)
		}
	})

	t.Run("ListByCapacity", func(t *testing.T) {
		auds, err := mgr.ListByCapacity(ctx, 20)
		assert.NoError(t, err)
		assert.NotEmpty(t, auds)

		found := false
		for _, a := range auds {
			if a.AudienceID == testAudience.AudienceID {
				found = true
				assert.GreaterOrEqual(t, a.Capacity, 20)
			}
		}
		assert.True(t, found, "created audience should be in ListByCapacity result")
	})

	t.Run("CheckNumberUnique", func(t *testing.T) {
		unique, err := mgr.CheckNumberUnique(ctx, "NonExistingNumber", 0)
		assert.NoError(t, err)
		assert.True(t, unique)

		unique, err = mgr.CheckNumberUnique(ctx, testAudience.AudinNumber, testAudience.AudienceID)
		assert.NoError(t, err)
		assert.True(t, unique)

		unique, err = mgr.CheckNumberUnique(ctx, testAudience.AudinNumber, 0)
		assert.NoError(t, err)
		assert.False(t, unique)
	})

	t.Run("Update", func(t *testing.T) {
		updatedAudience := *testAudience
		updatedAudience.Capacity = 60
		err := mgr.Update(ctx, &updatedAudience)
		assert.NoError(t, err)

		aud, err := mgr.GetByNumber(ctx, testAudience.AudinNumber)
		assert.NoError(t, err)
		if aud != nil {
			assert.Equal(t, 60, aud.Capacity)
		}
	})

	t.Run("Create duplicate number", func(t *testing.T) {
		dupAudience := &domain.Audience{
			Name:        "Test_duplicate",
			AudinType:   "Test_dupl",
			AudinNumber: testAudience.AudinNumber, // дублируем номер
			Capacity:    10,
		}
		err := mgr.Create(ctx, dupAudience)
		assert.Error(t, err)
	})
	if !t.Failed() {
		levelLogger.Info("All tests passed successfully")
		t.Log("All tests passed successfully")
	}
}

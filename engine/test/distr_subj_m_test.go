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

func TestSubjectDistributionManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewSubjectDistributionRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := managers.NewSubjectDistributionManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testDistribution := &domain.SubjectDistribution{
		EmployeeID: 5,
		SubjectID:  2,
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testDistribution)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		if testDistribution.SubjectDistrID == 0 {
			levelLogger.Error("Expected ID to be set after Create")
			t.Error("Expected ID to be set after Create")
		}
		levelLogger.Info("Created distribution", logger.Int("ID", testDistribution.SubjectDistrID))
	})

	t.Run("GetByEmployeeAndSubject", func(t *testing.T) {
		distr, err := mgr.GetByEmployeeAndSubject(ctx, testDistribution.EmployeeID, testDistribution.SubjectID)
		assert.NoError(t, err)
		assert.NotNil(t, distr)
		if distr != nil {
			assert.Equal(t, testDistribution.EmployeeID, distr.EmployeeID)
			assert.Equal(t, testDistribution.SubjectID, distr.SubjectID)
		}
	})

	t.Run("CheckExists", func(t *testing.T) {
		exists, err := mgr.CheckExists(ctx, testDistribution.EmployeeID, testDistribution.SubjectID)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = mgr.CheckExists(ctx, 9999, 9999) // Неверные ID
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("GetByEmployee", func(t *testing.T) {
		distributions, err := mgr.GetByEmployee(ctx, testDistribution.EmployeeID)
		assert.NoError(t, err)
		assert.NotEmpty(t, distributions)

		found := false
		for _, d := range distributions {
			if d.SubjectDistrID == testDistribution.SubjectDistrID {
				found = true
			}
		}
		assert.True(t, found, "created distribution should be in GetByEmployee result")
	})

	t.Run("GetBySubject", func(t *testing.T) {
		distributions, err := mgr.GetBySubject(ctx, testDistribution.SubjectID)
		assert.NoError(t, err)
		assert.NotEmpty(t, distributions)

		found := false
		for _, d := range distributions {
			if d.SubjectDistrID == testDistribution.SubjectDistrID {
				found = true
			}
		}
		assert.True(t, found, "created distribution should be in GetBySubject result")
	})

	t.Run("Create Duplicate", func(t *testing.T) {
		dupDistribution := &domain.SubjectDistribution{
			EmployeeID: testDistribution.EmployeeID,
			SubjectID:  testDistribution.SubjectID,
		}
		err := mgr.Create(ctx, dupDistribution)
		assert.Error(t, err)
	})

	if !t.Failed() {
		levelLogger.Info("All tests passed successfully")
		t.Log("All tests passed successfully")
	}
}

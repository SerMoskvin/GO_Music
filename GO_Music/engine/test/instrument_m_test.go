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

func TestInstrumentManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewInstrumentRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := managers.NewInstrumentManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testInstrument := &domain.Instrument{
		Name:       "TestInstrument1",
		InstrType:  "String",
		AudienceID: 1,
		Condition:  "Good",
		// Добавьте обязательные поля, если есть
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testInstrument)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		assert.NotZero(t, testInstrument.InstrumentID)
		levelLogger.Info("Created instrument", logger.Int("ID", testInstrument.InstrumentID))
	})

	t.Run("GetByName", func(t *testing.T) {
		instr, err := mgr.GetByName(ctx, testInstrument.Name)
		assert.NoError(t, err)
		assert.NotNil(t, instr)
		if instr != nil {
			assert.Equal(t, testInstrument.Name, instr.Name)
		}
	})

	t.Run("GetByType", func(t *testing.T) {
		instruments, err := mgr.GetByType(ctx, testInstrument.InstrType)
		assert.NoError(t, err)
		assert.NotEmpty(t, instruments)

		found := false
		for _, i := range instruments {
			if i.InstrumentID == testInstrument.InstrumentID {
				found = true
			}
		}
		assert.True(t, found, "Created instrument should be in GetByType result")
	})

	t.Run("GetByAudience", func(t *testing.T) {
		instruments, err := mgr.GetByAudience(ctx, testInstrument.AudienceID)
		assert.NoError(t, err)
		assert.NotEmpty(t, instruments)

		found := false
		for _, i := range instruments {
			if i.InstrumentID == testInstrument.InstrumentID {
				found = true
			}
		}
		assert.True(t, found, "Created instrument should be in GetByAudience result")
	})

	t.Run("CheckNameUnique", func(t *testing.T) {
		isUnique, err := mgr.CheckNameUnique(ctx, testInstrument.Name, testInstrument.InstrumentID)
		assert.NoError(t, err)
		assert.True(t, isUnique, "Name should be unique excluding current instrument")

		isUnique, err = mgr.CheckNameUnique(ctx, testInstrument.Name, 0)
		assert.NoError(t, err)
		assert.False(t, isUnique, "Name should not be unique without exclusion")

		isUnique, err = mgr.CheckNameUnique(ctx, "NonexistentInstrumentName", 0)
		assert.NoError(t, err)
		assert.True(t, isUnique, "Nonexistent instrument name should be unique")
	})

	t.Run("UpdateCondition", func(t *testing.T) {
		newCondition := "Excellent"
		err := mgr.UpdateCondition(ctx, testInstrument.InstrumentID, newCondition)
		assert.NoError(t, err)

		instr, err := mgr.GetByID(ctx, testInstrument.InstrumentID)
		assert.NoError(t, err)
		assert.NotNil(t, instr)
		if instr != nil {
			assert.Equal(t, newCondition, instr.Condition)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newName := "UpdatedInstrumentName"
		testInstrument.Name = newName
		err := mgr.Update(ctx, testInstrument)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		instr, err := mgr.GetByName(ctx, newName)
		assert.NoError(t, err)
		assert.NotNil(t, instr)
		if instr != nil {
			assert.Equal(t, newName, instr.Name)
		}
	})

	t.Run("Create Duplicate Name", func(t *testing.T) {
		dupInstrument := &domain.Instrument{
			Name:       testInstrument.Name,
			InstrType:  "Wind",
			AudienceID: 2,
			Condition:  "Fair",
		}
		err := mgr.Create(ctx, dupInstrument)
		assert.Error(t, err)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		instruments := []*domain.Instrument{
			{
				Name:       "BulkInstrument1",
				InstrType:  "Percussion",
				AudienceID: 1,
				Condition:  "Good",
			},
			{
				Name:       "BulkInstrument2",
				InstrType:  "Keyboard",
				AudienceID: 1,
				Condition:  "Good",
			},
		}

		err := mgr.BulkCreate(ctx, instruments)
		assert.NoError(t, err)

		for _, i := range instruments {
			assert.NotZero(t, i.InstrumentID)
		}
	})

	if !t.Failed() {
		levelLogger.Info("All InstrumentManager tests passed successfully")
		t.Log("All InstrumentManager tests passed successfully")
	}
}

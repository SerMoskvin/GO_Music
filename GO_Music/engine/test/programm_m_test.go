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

func TestProgrammManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewProgrammRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := engine.NewProgrammManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testProgramm := &domain.Programm{
		ProgrammName:           "TestProgramm_",
		ProgrammType:           "Classical",
		Duration:               60,
		Instrument:             ptrString("Piano"),
		Description:            ptrString("Test description"),
		StudyLoad:              10,
		FinalCertificationForm: "Exam",
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testProgramm)
		assert.NoError(t, err)
		assert.NotZero(t, testProgramm.MusprogrammID)
	})

	t.Run("GetByType", func(t *testing.T) {
		progs, err := mgr.GetByType(ctx, testProgramm.ProgrammType)
		assert.NoError(t, err)
		assert.NotEmpty(t, progs)

		found := false
		for _, p := range progs {
			if p.MusprogrammID == testProgramm.MusprogrammID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created programm should be in GetByType result")
	})

	t.Run("GetByInstrument", func(t *testing.T) {
		progs, err := mgr.GetByInstrument(ctx, *testProgramm.Instrument)
		assert.NoError(t, err)
		assert.NotEmpty(t, progs)

		found := false
		for _, p := range progs {
			if p.MusprogrammID == testProgramm.MusprogrammID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created programm should be in GetByInstrument result")
	})

	t.Run("GetByName", func(t *testing.T) {
		p, err := mgr.GetByName(ctx, testProgramm.ProgrammName)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, testProgramm.MusprogrammID, p.MusprogrammID)
	})

	t.Run("GetByDurationRange", func(t *testing.T) {
		progs, err := mgr.GetByDurationRange(ctx, 30, 90)
		assert.NoError(t, err)
		assert.NotEmpty(t, progs)

		found := false
		for _, p := range progs {
			if p.MusprogrammID == testProgramm.MusprogrammID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("GetByStudyLoad", func(t *testing.T) {
		progs, err := mgr.GetByStudyLoad(ctx, testProgramm.StudyLoad)
		assert.NoError(t, err)
		assert.NotEmpty(t, progs)

		found := false
		for _, p := range progs {
			if p.MusprogrammID == testProgramm.MusprogrammID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("CheckNameUnique", func(t *testing.T) {
		unique, err := mgr.CheckNameUnique(ctx, testProgramm.ProgrammName, testProgramm.MusprogrammID)
		assert.NoError(t, err)
		assert.True(t, unique, "Name should be unique excluding self")

		unique, err = mgr.CheckNameUnique(ctx, testProgramm.ProgrammName, 0)
		assert.NoError(t, err)
		assert.False(t, unique, "Name should not be unique if excludeID=0")
	})

	t.Run("SearchByDescription", func(t *testing.T) {
		progs, err := mgr.SearchByDescription(ctx, "description")
		assert.NoError(t, err)
		assert.NotEmpty(t, progs)

		found := false
		for _, p := range progs {
			if p.MusprogrammID == testProgramm.MusprogrammID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Update", func(t *testing.T) {
		testProgramm.Description = ptrString("Updated description")
		err := mgr.Update(ctx, testProgramm)
		assert.NoError(t, err)

		p, err := mgr.GetByName(ctx, testProgramm.ProgrammName)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", *p.Description)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		progs := []*domain.Programm{
			{
				ProgrammName:           "BulkProg1_",
				ProgrammType:           "Jazz",
				Duration:               45,
				Instrument:             ptrString("Saxophone"),
				Description:            ptrString("Bulk description 1"),
				StudyLoad:              5,
				FinalCertificationForm: "Test",
			},
			{
				ProgrammName:           "BulkProg2_",
				ProgrammType:           "Rock",
				Duration:               50,
				Instrument:             ptrString("Guitar"),
				Description:            ptrString("Bulk description 2"),
				StudyLoad:              6,
				FinalCertificationForm: "Test",
			},
		}

		err := mgr.BulkCreate(ctx, progs)
		assert.NoError(t, err)

		for _, p := range progs {
			assert.NotZero(t, p.MusprogrammID)
		}
	})

	if !t.Failed() {
		levelLogger.Info("All ProgrammManager tests passed successfully")
		t.Log("All ProgrammManager tests passed successfully")
	}
}

func ptrString(s string) *string {
	return &s
}

func ptrNumber(n int) *int {
	return &n
}

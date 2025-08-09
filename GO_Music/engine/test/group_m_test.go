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

func TestStudyGroupManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudyGroupRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := engine.NewStudyGroupManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testGroup := &domain.StudyGroup{
		GroupName:        "TestGroup1",
		MusProgrammID:    1,
		StudyYear:        2023,
		NumberOfStudents: 10,
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testGroup)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		if testGroup.GroupID == 0 {
			levelLogger.Error("Expected ID to be set after Create")
			t.Error("Expected ID to be set after Create")
		}
		levelLogger.Info("Created study group", logger.Int("ID", testGroup.GroupID))
	})

	t.Run("GetByName", func(t *testing.T) {
		gr, err := mgr.GetByName(ctx, testGroup.GroupName)
		assert.NoError(t, err)
		assert.NotNil(t, gr)
		if gr != nil {
			assert.Equal(t, testGroup.GroupName, gr.GroupName)
		}
	})

	t.Run("GetByProgram", func(t *testing.T) {
		groups, err := mgr.GetByProgram(ctx, testGroup.MusProgrammID)
		assert.NoError(t, err)
		assert.NotEmpty(t, groups)

		found := false
		for _, g := range groups {
			if g.GroupID == testGroup.GroupID {
				found = true
			}
		}
		assert.True(t, found, "Created group should be in GetByProgram result")
	})

	t.Run("GetByYear", func(t *testing.T) {
		groups, err := mgr.GetByYear(ctx, testGroup.StudyYear)
		assert.NoError(t, err)
		assert.NotEmpty(t, groups)

		found := false
		for _, g := range groups {
			if g.GroupID == testGroup.GroupID {
				found = true
			}
		}
		assert.True(t, found, "Created group should be in GetByYear result")
	})

	t.Run("CheckNameUnique", func(t *testing.T) {
		isUnique, err := mgr.CheckNameUnique(ctx, testGroup.GroupName, testGroup.GroupID)
		assert.NoError(t, err)
		assert.True(t, isUnique, "Name should be unique excluding current group")

		isUnique, err = mgr.CheckNameUnique(ctx, testGroup.GroupName, 0)
		assert.NoError(t, err)
		assert.False(t, isUnique, "Name should not be unique without exclusion")

		isUnique, err = mgr.CheckNameUnique(ctx, "NonexistentGroupName", 0)
		assert.NoError(t, err)
		assert.True(t, isUnique, "Nonexistent group name should be unique")
	})

	t.Run("UpdateStudentCount", func(t *testing.T) {
		newCount := testGroup.NumberOfStudents + 5
		err := mgr.UpdateStudentCount(ctx, testGroup.GroupID, newCount)
		assert.NoError(t, err)

		gr, err := mgr.GetByID(ctx, testGroup.GroupID)
		assert.NoError(t, err)
		assert.NotNil(t, gr)
		if gr != nil {
			assert.Equal(t, newCount, gr.NumberOfStudents)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newName := "UpdatedGroupName"
		testGroup.GroupName = newName
		err := mgr.Update(ctx, testGroup)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		gr, err := mgr.GetByName(ctx, newName)
		assert.NoError(t, err)
		assert.NotNil(t, gr)
		if gr != nil {
			assert.Equal(t, newName, gr.GroupName)
		}
	})

	t.Run("Create Duplicate Name", func(t *testing.T) {
		dupGroup := &domain.StudyGroup{
			GroupName:        testGroup.GroupName,
			MusProgrammID:    2,
			StudyYear:        2024,
			NumberOfStudents: 5,
		}
		err := mgr.Create(ctx, dupGroup)
		assert.Error(t, err)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		groups := []*domain.StudyGroup{
			{
				GroupName:        "BulkGroup1",
				MusProgrammID:    1,
				StudyYear:        2022,
				NumberOfStudents: 15,
			},
			{
				GroupName:        "BulkGroup2",
				MusProgrammID:    1,
				StudyYear:        2023,
				NumberOfStudents: 20,
			},
		}

		err := mgr.BulkCreate(ctx, groups)
		assert.NoError(t, err)

		for _, g := range groups {
			assert.NotZero(t, g.GroupID)
		}
	})

	if !t.Failed() {
		levelLogger.Info("All StudyGroupManager tests passed successfully")
		t.Log("All StudyGroupManager tests passed successfully")
	}
}

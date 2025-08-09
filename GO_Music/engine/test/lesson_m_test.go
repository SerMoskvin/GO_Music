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

func TestLessonManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewLessonRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := managers.NewLessonManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	StudentID, AudienceID := 4, 3
	testLesson := &domain.Lesson{
		LessonName: "Testlesson",
		EmployeeID: 7,
		GroupID:    1,
		StudentID:  &StudentID,
		SubjectID:  1,
		AudienceID: &AudienceID,
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.BulkCreate(ctx, []*domain.Lesson{testLesson})
		assert.NoError(t, err)
		assert.NotZero(t, testLesson.LessonID)
	})

	t.Run("GetByEmployee", func(t *testing.T) {
		lessons, err := mgr.GetByEmployee(ctx, testLesson.EmployeeID)
		assert.NoError(t, err)
		assert.NotEmpty(t, lessons)

		found := false
		for _, l := range lessons {
			if l.LessonID == testLesson.LessonID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created lesson should be in GetByEmployee result")
	})

	t.Run("GetByGroup", func(t *testing.T) {
		lessons, err := mgr.GetByGroup(ctx, testLesson.GroupID)
		assert.NoError(t, err)
		assert.NotEmpty(t, lessons)

		found := false
		for _, l := range lessons {
			if l.LessonID == testLesson.LessonID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created lesson should be in GetByGroup result")
	})

	t.Run("GetByStudent", func(t *testing.T) {
		lessons, err := mgr.GetByStudent(ctx, *testLesson.StudentID)
		assert.NoError(t, err)
		assert.NotEmpty(t, lessons)

		found := false
		for _, l := range lessons {
			if l.LessonID == testLesson.LessonID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created lesson should be in GetByStudent result")
	})

	t.Run("GetBySubject", func(t *testing.T) {
		lessons, err := mgr.GetBySubject(ctx, testLesson.SubjectID)
		assert.NoError(t, err)
		assert.NotEmpty(t, lessons)

		found := false
		for _, l := range lessons {
			if l.LessonID == testLesson.LessonID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created lesson should be in GetBySubject result")
	})

	t.Run("GetByAudience", func(t *testing.T) {
		lessons, err := mgr.GetByAudience(ctx, *testLesson.AudienceID)
		assert.NoError(t, err)
		assert.NotEmpty(t, lessons)

		found := false
		for _, l := range lessons {
			if l.LessonID == testLesson.LessonID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created lesson should be in GetByAudience result")
	})

	now := time.Now()
	start := now
	end := now.Add(1 * time.Hour)

	t.Run("CheckEmployeeAvailability", func(t *testing.T) {
		available, err := mgr.CheckEmployeeAvailability(ctx, testLesson.EmployeeID, start, end, 0)
		if err != nil {
			t.Errorf("CheckEmployeeAvailability returned error: %v", err)
			levelLogger.Error("CheckEmployeeAvailability error", logger.Field{Key: "error", Value: err})
			return
		}
		t.Logf("Employee availability: %v", available)
	})

	t.Run("CheckAudienceAvailability", func(t *testing.T) {
		available, err := mgr.CheckAudienceAvailability(ctx, *testLesson.AudienceID, start, end, 0)
		if err != nil {
			t.Errorf("CheckAudienceAvailability returned error: %v", err)
			levelLogger.Error("CheckAudienceAvailability error", logger.Field{Key: "error", Value: err})
			return
		}
		t.Logf("Audience availability: %v", available)
	})

	t.Run("BulkCreate", func(t *testing.T) {
		lessons := []*domain.Lesson{
			{
				LessonName: "BulkLesson1",
				EmployeeID: 6,
				GroupID:    3,
				StudentID:  &StudentID,
				SubjectID:  3,
				AudienceID: &AudienceID,
			},
			{
				LessonName: "BulkLesson2",
				EmployeeID: 7,
				GroupID:    2,
				StudentID:  &StudentID,
				SubjectID:  1,
				AudienceID: &AudienceID,
			},
		}

		err := mgr.BulkCreate(ctx, lessons)
		assert.NoError(t, err)

		for _, l := range lessons {
			assert.NotZero(t, l.LessonID)
		}
	})

	if !t.Failed() {
		levelLogger.Info("All LessonManager tests passed successfully")
		t.Log("All LessonManager tests passed successfully")
	}
}

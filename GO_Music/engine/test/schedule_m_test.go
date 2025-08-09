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

func TestScheduleManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewScheduleRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := engine.NewScheduleManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testSchedule := &domain.Schedule{
		LessonID:      3,
		DayWeek:       "Понедельник",
		TimeBegin:     domain.ParseTimeHM("10:00"),
		TimeEnd:       domain.ParseTimeHM("11:00"),
		SchdDateStart: domain.ParseDMY("30.10.2023"),
		SchdDateEnd:   domain.ParseDMY("30.11.2023"),
	}

	t.Run("Create", func(t *testing.T) {
		err := mgr.Create(ctx, testSchedule)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}
		assert.NotZero(t, testSchedule.ScheduleID)
	})

	t.Run("GetByLesson", func(t *testing.T) {
		schedules, err := mgr.GetByLesson(ctx, testSchedule.LessonID)
		if err != nil {
			levelLogger.Error("GetByLesson failed", logger.String("error", err.Error()))
			t.Fatalf("GetByLesson failed: %v", err)
		}
		assert.NotEmpty(t, schedules)

		found := false
		for _, s := range schedules {
			if s.ScheduleID == testSchedule.ScheduleID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created schedule should be in GetByLesson result")
	})

	t.Run("GetByDay", func(t *testing.T) {
		schedules, err := mgr.GetByDay(ctx, testSchedule.DayWeek)
		if err != nil {
			levelLogger.Error("GetByDay failed", logger.String("error", err.Error()))
			t.Fatalf("GetByDay failed: %v", err)
		}
		assert.NotEmpty(t, schedules)

		found := false
		for _, s := range schedules {
			if s.ScheduleID == testSchedule.ScheduleID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created schedule should be in GetByDay result")
	})

	t.Run("GetCurrentSchedule", func(t *testing.T) {
		schedules, err := mgr.GetCurrentSchedule(ctx)
		if err != nil {
			levelLogger.Error("GetCurrentSchedule failed", logger.String("error", err.Error()))
			t.Fatalf("GetCurrentSchedule failed: %v", err)
		}
		assert.NotEmpty(t, schedules)
	})

	t.Run("CheckTimeConflict", func(t *testing.T) {
		conflict, err := mgr.CheckTimeConflict(ctx, testSchedule.DayWeek, testSchedule.TimeBegin.Format("15:04"), testSchedule.TimeEnd.Format("15:04"), testSchedule.ScheduleID)
		if err != nil {
			levelLogger.Error("CheckTimeConflict failed", logger.String("error", err.Error()))
			t.Fatalf("CheckTimeConflict failed: %v", err)
		}
		assert.False(t, conflict, "There should be no time conflict for the created schedule")
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		startDate := domain.ParseDMY("01.08.2023")
		endDate := domain.ParseDMY("01.09.2025")
		schedules, err := mgr.GetByDateRange(ctx, startDate, endDate)
		if err != nil {
			levelLogger.Error("GetByDateRange failed", logger.String("error", err.Error()))
			t.Fatalf("GetByDateRange failed: %v", err)
		}
		assert.NotEmpty(t, schedules)

		found := false
		for _, s := range schedules {
			if s.ScheduleID == testSchedule.ScheduleID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Update", func(t *testing.T) {
		testSchedule.TimeEnd = testSchedule.TimeEnd.Add(30 * time.Minute)
		err := mgr.Update(ctx, testSchedule)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		updatedSchedule, err := mgr.GetByLesson(ctx, testSchedule.LessonID)
		if err != nil {
			levelLogger.Error("GetByLesson after update failed", logger.String("error", err.Error()))
			t.Fatalf("GetByLesson after update failed: %v", err)
		}
		assert.Equal(t, testSchedule.TimeEnd, updatedSchedule[0].TimeEnd)
	})

	if !t.Failed() {
		levelLogger.Info("All ScheduleManager tests passed successfully")
		t.Log("All ScheduleManager tests passed successfully")
	}
}

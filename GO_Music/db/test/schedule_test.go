package db_test

import (
	"context"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	_ "github.com/lib/pq"
)

func TestScheduleRepository_AllMethods(t *testing.T) {
	cfgPath := "../../config/DB_config.yml"
	cfg, err := config.LoadDBConfig(cfgPath)
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

	// Создаем тестовые даты
	startDate := time.Now().Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 1, 0) // +1 месяц

	testSchedule := &domain.Schedule{
		LessonID:      3,
		DayWeek:       "Monday",
		TimeBegin:     domain.ParseTimeHM("09:00"),
		TimeEnd:       domain.ParseTimeHM("10:30"),
		SchdDateStart: startDate,
		SchdDateEnd:   endDate,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testSchedule)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testSchedule.ScheduleID == 0 {
			t.Error("Expected ScheduleID to be set after Create")
		} else {
			t.Logf("Created ScheduleID: %d", testSchedule.ScheduleID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		sched, err := repo.GetByID(ctx, testSchedule.ScheduleID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if sched == nil {
			t.Fatal("Expected Schedule to be found")
		}
		if sched.LessonID != testSchedule.LessonID {
			t.Errorf("Expected LessonID %d, got %d", testSchedule.LessonID, sched.LessonID)
		}
		if sched.DayWeek != testSchedule.DayWeek {
			t.Errorf("Expected DayWeek %q, got %q", testSchedule.DayWeek, sched.DayWeek)
		}
		if sched.TimeBegin != testSchedule.TimeBegin {
			t.Errorf("Expected TimeBegin %q, got %q", testSchedule.TimeBegin, sched.TimeBegin)
		}
		if sched.TimeEnd != testSchedule.TimeEnd {
			t.Errorf("Expected TimeEnd %q, got %q", testSchedule.TimeEnd, sched.TimeEnd)
		}
		if !sched.SchdDateStart.Equal(testSchedule.SchdDateStart) {
			t.Errorf("Expected SchedDateStart %v, got %v", testSchedule.SchdDateStart, sched.SchdDateStart)
		}
		if !sched.SchdDateEnd.Equal(testSchedule.SchdDateEnd) {
			t.Errorf("Expected SchedDateEnd %v, got %v", testSchedule.SchdDateEnd, sched.SchdDateEnd)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedSchedule := *testSchedule
		updatedSchedule.DayWeek = "Tuesday"
		updatedSchedule.TimeBegin = domain.ParseTimeHM("10:00")
		updatedSchedule.TimeEnd = domain.ParseTimeHM("11:00")
		updatedSchedule.SchdDateStart = updatedSchedule.SchdDateStart.AddDate(0, 0, 1) // +1 день
		updatedSchedule.SchdDateEnd = updatedSchedule.SchdDateEnd.AddDate(0, 0, 1)     // +1 день

		err = repoWithTx.Update(ctx, &updatedSchedule)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		sched, err := repo.GetByID(ctx, testSchedule.ScheduleID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if sched.DayWeek != updatedSchedule.DayWeek {
			t.Errorf("Expected DayWeek %q after update, got %q", updatedSchedule.DayWeek, sched.DayWeek)
		}
		if sched.TimeBegin != updatedSchedule.TimeBegin {
			t.Errorf("Expected TimeBegin %q after update, got %q", updatedSchedule.TimeBegin, sched.TimeBegin)
		}
		if sched.TimeEnd != updatedSchedule.TimeEnd {
			t.Errorf("Expected TimeEnd %q after update, got %q", updatedSchedule.TimeEnd, sched.TimeEnd)
		}
		if !sched.SchdDateStart.Equal(updatedSchedule.SchdDateStart) {
			t.Errorf("Expected SchedDateStart %v after update, got %v", updatedSchedule.SchdDateStart, sched.SchdDateStart)
		}
		if !sched.SchdDateEnd.Equal(updatedSchedule.SchdDateEnd) {
			t.Errorf("Expected SchedDateEnd %v after update, got %v", updatedSchedule.SchdDateEnd, sched.SchdDateEnd)
		}

		*testSchedule = updatedSchedule
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "day_week",
					Operator: "=",
					Value:    testSchedule.DayWeek,
				},
			},
			Limit: 10,
		}

		schedules, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(schedules) == 0 {
			t.Error("Expected at least one Schedule in List")
		} else {
			t.Logf("List returned %d items", len(schedules))
		}

		found := false
		for _, s := range schedules {
			if s.ScheduleID == testSchedule.ScheduleID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Schedule not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "day_week",
					Operator: "=",
					Value:    testSchedule.DayWeek,
				},
			},
		}

		count, err := repo.Count(ctx, filter)
		if err != nil {
			t.Fatalf("Count failed: %v", err)
		}

		if count < 1 {
			t.Errorf("Expected count >= 1, got %d", count)
		} else {
			t.Logf("Count returned %d", count)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, testSchedule.ScheduleID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Schedule to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Schedule to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondSchedule := &domain.Schedule{
			LessonID:      2, // валидный LessonID
			DayWeek:       "Wednesday",
			TimeBegin:     domain.ParseTimeHM("11:00"),
			TimeEnd:       domain.ParseTimeHM("12:00"),
			SchdDateStart: startDate,
			SchdDateEnd:   endDate,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondSchedule)
		if err != nil {
			t.Fatalf("Create second Schedule failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testSchedule.ScheduleID, secondSchedule.ScheduleID}
		schedules, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(schedules) != 2 {
			t.Errorf("Expected 2 Schedules, got %d", len(schedules))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testSchedule.ScheduleID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testSchedule.ScheduleID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Schedule to be deleted")
		}
	})
}

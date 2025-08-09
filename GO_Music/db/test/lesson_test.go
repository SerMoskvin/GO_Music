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

func TestLessonRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewLessonRepository(sqlDB)

	// укажите валидные IDs для FK
	audienceID := 1
	employeeID := 6
	groupID := 1
	subjectID := 1

	testLesson := &domain.Lesson{
		AudienceID: &audienceID,
		EmployeeID: employeeID,
		GroupID:    groupID,
		StudentID:  nil,
		LessonName: "Урок музыки",
		SubjectID:  subjectID,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testLesson)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testLesson.LessonID == 0 {
			t.Error("Expected LessonID to be set after Create")
		} else {
			t.Logf("Created LessonID: %d", testLesson.LessonID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		lesson, err := repo.GetByID(ctx, testLesson.LessonID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if lesson == nil {
			t.Fatal("Expected Lesson to be found")
		}
		if lesson.LessonName != testLesson.LessonName {
			t.Errorf("Expected LessonName %q, got %q", testLesson.LessonName, lesson.LessonName)
		}
		if lesson.EmployeeID != testLesson.EmployeeID {
			t.Errorf("Expected EmployeeID %d, got %d", testLesson.EmployeeID, lesson.EmployeeID)
		}
		if lesson.GroupID != testLesson.GroupID {
			t.Errorf("Expected GroupID %d, got %d", testLesson.GroupID, lesson.GroupID)
		}
		if lesson.SubjectID != testLesson.SubjectID {
			t.Errorf("Expected SubjectID %d, got %d", testLesson.SubjectID, lesson.SubjectID)
		}
		if testLesson.AudienceID != nil && lesson.AudienceID != nil && *lesson.AudienceID != *testLesson.AudienceID {
			t.Errorf("Expected AudienceID %d, got %d", *testLesson.AudienceID, *lesson.AudienceID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedLesson := *testLesson
		updatedLesson.LessonName = "Обновленный урок"

		err = repoWithTx.Update(ctx, &updatedLesson)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		lesson, err := repo.GetByID(ctx, testLesson.LessonID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if lesson.LessonName != "Обновленный урок" {
			t.Errorf("Expected LessonName 'Обновленный урок' after update, got %q", lesson.LessonName)
		}

		*testLesson = updatedLesson
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "lesson_name",
					Operator: "=",
					Value:    testLesson.LessonName,
				},
			},
			Limit: 10,
		}

		lessons, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(lessons) == 0 {
			t.Error("Expected at least one Lesson in List")
		} else {
			t.Logf("List returned %d items", len(lessons))
		}

		found := false
		for _, l := range lessons {
			if l.LessonID == testLesson.LessonID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Lesson not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "lesson_name",
					Operator: "=",
					Value:    testLesson.LessonName,
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
		exists, err := repo.Exists(ctx, testLesson.LessonID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Lesson to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Lesson to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondLesson := &domain.Lesson{
			AudienceID: &audienceID,
			EmployeeID: employeeID,
			GroupID:    groupID,
			StudentID:  nil,
			LessonName: "Второй урок",
			SubjectID:  subjectID,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondLesson)
		if err != nil {
			t.Fatalf("Create second Lesson failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testLesson.LessonID, secondLesson.LessonID}
		lessons, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(lessons) != 2 {
			t.Errorf("Expected 2 Lessons, got %d", len(lessons))
		}
	})

	t.Run("CheckEmployeeAvailability", func(t *testing.T) {
		// Для корректного теста нужно, чтобы в таблицах lessons и schedules были данные.
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)

		available, err := repo.CheckEmployeeAvailability(ctx, employeeID, startTime, endTime, 0)
		if err != nil {
			t.Fatalf("CheckEmployeeAvailability failed: %v", err)
		}

		t.Logf("Employee availability: %v", available)
	})

	t.Run("CheckAudienceAvailability", func(t *testing.T) {
		if testLesson.AudienceID == nil {
			t.Skip("AudienceID is nil, skipping CheckAudienceAvailability test")
		}
		startTime := time.Now().Add(1 * time.Hour)
		endTime := startTime.Add(1 * time.Hour)

		available, err := repo.CheckAudienceAvailability(ctx, *testLesson.AudienceID, startTime, endTime, 0)
		if err != nil {
			t.Fatalf("CheckAudienceAvailability failed: %v", err)
		}

		t.Logf("Audience availability: %v", available)
	})
}

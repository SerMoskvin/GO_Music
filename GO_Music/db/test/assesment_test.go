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

func TestStudentAssessmentRepository_AllMethods(t *testing.T) {
	// Загрузка конфигурации
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

	repo := repositories.NewStudentAssessmentRepository(sqlDB)

	// Тестовые данные
	testAssessment := &domain.StudentAssessment{
		StudentID:      2,
		LessonID:       2,
		Grade:          5,
		TaskType:       "test",
		AssessmentDate: time.Now(),
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testAssessment)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testAssessment.AssessmentNoteID == 0 {
			t.Error("Expected AssessmentNoteID to be set after Create")
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		assessment, err := repo.GetByID(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if assessment == nil {
			t.Fatal("Expected assessment to be found")
		}

		if assessment.StudentID != testAssessment.StudentID {
			t.Errorf("Expected StudentID %d, got %d", testAssessment.StudentID, assessment.StudentID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedAssessment := *testAssessment
		updatedAssessment.Grade = 4
		updatedAssessment.TaskType = "updated test"

		err = repoWithTx.Update(ctx, &updatedAssessment)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		assessment, err := repo.GetByID(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if assessment.Grade != 4 {
			t.Errorf("Expected Grade 4 after update, got %d", assessment.Grade)
		}
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "student_id",
					Operator: "=",
					Value:    testAssessment.StudentID,
				},
			},
			Limit: 10,
		}

		assessments, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(assessments) == 0 {
			t.Error("Expected at least one assessment in List")
		}

		found := false
		for _, a := range assessments {
			if a.AssessmentNoteID == testAssessment.AssessmentNoteID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Created assessment not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "student_id",
					Operator: "=",
					Value:    testAssessment.StudentID,
				},
			},
		}

		count, err := repo.Count(ctx, filter)
		if err != nil {
			t.Fatalf("Count failed: %v", err)
		}

		if count < 1 {
			t.Errorf("Expected count >= 1, got %d", count)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}

		if !exists {
			t.Error("Expected assessment to exist")
		}

		// Проверка несуществующей записи
		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}

		if exists {
			t.Error("Expected non-existent assessment to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondAssessment := &domain.StudentAssessment{
			StudentID:      3,
			LessonID:       3,
			Grade:          4,
			TaskType:       "test_2",
			AssessmentDate: time.Now(),
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondAssessment)
		if err != nil {
			t.Fatalf("Create second assessment failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testAssessment.AssessmentNoteID, secondAssessment.AssessmentNoteID}
		assessments, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(assessments) != 2 {
			t.Errorf("Expected 2 assessments, got %d", len(assessments))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			t.Error("Expected assessment to be deleted")
		}
	})
}

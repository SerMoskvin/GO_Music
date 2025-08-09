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
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestStudentAssessmentManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentAssessmentRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := engine.NewStudentAssessmentManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testAssessment := &domain.StudentAssessment{
		StudentID:      2,
		LessonID:       2,
		Grade:          5,
		TaskType:       "test",
		AssessmentDate: time.Now(),
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, testAssessment)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}

		if testAssessment.AssessmentNoteID == 0 {
			levelLogger.Error("Expected AssessmentNoteID to be set after Create")
			t.Error("Expected AssessmentNoteID to be set after Create")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		assessment, err := repo.GetByID(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			levelLogger.Error("GetByID failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID failed: %v", err)
		}

		if assessment == nil {
			levelLogger.Error("Expected assessment to be found")
			t.Fatal("Expected assessment to be found")
		}

		if assessment.StudentID != testAssessment.StudentID {
			levelLogger.Error("StudentID mismatch", logger.Int("expected", testAssessment.StudentID), logger.Int("got", assessment.StudentID))
			t.Errorf("Expected StudentID %d, got %d", testAssessment.StudentID, assessment.StudentID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		updatedAssessment := *testAssessment
		updatedAssessment.Grade = 4
		updatedAssessment.TaskType = "updated test"

		err = repo.Update(ctx, &updatedAssessment)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		assessment, err := repo.GetByID(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			levelLogger.Error("GetByID after update failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if assessment.Grade != 4 {
			levelLogger.Error("Grade mismatch after update", logger.Int("expected", 4), logger.Int("got", assessment.Grade))
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
			levelLogger.Error("List failed", logger.String("error", err.Error()))
			t.Fatalf("List failed: %v", err)
		}

		if len(assessments) == 0 {
			levelLogger.Error("Expected at least one assessment in List")
			t.Error("Expected at least one assessment in List")
		}
	})

	t.Run("BulkUpsert", func(t *testing.T) {
		assessments := []*domain.StudentAssessment{
			{StudentID: 2, LessonID: 2, Grade: 5, TaskType: "test", AssessmentDate: time.Now()},
			{StudentID: 2, LessonID: 3, Grade: 4, TaskType: "test", AssessmentDate: time.Now()},
		}

		err := mgr.BulkUpsert(ctx, assessments)
		if err != nil {
			levelLogger.Error("BulkUpsert failed", logger.String("error", err.Error()))
		}
		assert.NoError(t, err)

		for _, assessment := range assessments {
			exists, err := repo.Exists(ctx, assessment.GetID())
			if err != nil {
				levelLogger.Error("Exists check failed after BulkUpsert", logger.String("error", err.Error()), logger.Int("AssessmentNoteID", assessment.GetID()))
			}
			assert.NoError(t, err)

			if !exists {
				levelLogger.Error("Expected assessment to exist after BulkUpsert", logger.Int("AssessmentNoteID", assessment.GetID()))
			}
			assert.True(t, exists, "Expected assessment to exist after BulkUpsert")
		}
	})

	t.Run("GetByLesson", func(t *testing.T) {
		lessonID := testAssessment.LessonID
		assessments, err := mgr.GetByLesson(ctx, lessonID)
		if err != nil {
			levelLogger.Error("GetByLesson failed", logger.String("error", err.Error()), logger.Int("lessonID", lessonID))
		}
		assert.NoError(t, err)

		if len(assessments) == 0 {
			levelLogger.Error("Expected assessments for lesson", logger.Int("lessonID", lessonID))
		}
		assert.NotEmpty(t, assessments, "Expected assessments for lesson")
	})

	t.Run("GetStudentAverageGrade", func(t *testing.T) {
		avg, err := mgr.GetStudentAverageGrade(ctx, testAssessment.StudentID)
		if err != nil {
			levelLogger.Error("GetStudentAverageGrade failed", logger.String("error", err.Error()), logger.Int("StudentID", testAssessment.StudentID))
		}
		assert.NoError(t, err)

		if avg <= 0.0 {
			levelLogger.Error("Expected average grade to be greater than 0", logger.Float64("averageGrade", avg), logger.Int("StudentID", testAssessment.StudentID))
		}
		assert.Greater(t, avg, 0.0, "Expected average grade to be greater than 0")
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			levelLogger.Error("Delete failed", logger.String("error", err.Error()))
			t.Fatalf("Delete failed: %v", err)
		}

		exists, err := repo.Exists(ctx, testAssessment.AssessmentNoteID)
		if err != nil {
			levelLogger.Error("Exists after delete failed", logger.String("error", err.Error()))
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			levelLogger.Error("Expected assessment to be deleted")
			t.Error("Expected assessment to be deleted")
		}
	})

	levelLogger.Info("All tests passed successfully")
	t.Log("All tests passed successfully")
}

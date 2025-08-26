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
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestStudentAttendanceManager_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentAttendanceRepository(sqlDB)

	levelLogger, err := logger.NewLevel(cfgPath_Log)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer levelLogger.Sync()

	mgr := managers.NewStudentAttendanceManager(repo, sqlDB, levelLogger, 5*time.Second)

	// Тестовые данные
	testAttendance := &domain.StudentAttendance{
		StudentID:      2,
		LessonID:       2,
		PresenceMark:   true,
		AttendanceDate: time.Now().Format("2006-01-02"),
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, testAttendance)
		if err != nil {
			levelLogger.Error("Create failed", logger.String("error", err.Error()))
			t.Fatalf("Create failed: %v", err)
		}

		if testAttendance.AttendanceNoteID == 0 {
			levelLogger.Error("Expected ID to be set after Create")
			t.Error("Expected ID to be set after Create")
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		attendance, err := repo.GetByID(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			levelLogger.Error("GetByID failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID failed: %v", err)
		}

		if attendance == nil {
			levelLogger.Error("Expected attendance to be found")
			t.Fatal("Expected attendance to be found")
		}

		if attendance.StudentID != testAttendance.StudentID {
			levelLogger.Error("StudentID mismatch", logger.Int("expected", testAttendance.StudentID), logger.Int("got", attendance.StudentID))
			t.Errorf("Expected StudentID %d, got %d", testAttendance.StudentID, attendance.StudentID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		updatedAttendance := *testAttendance
		updatedAttendance.PresenceMark = false

		err = repo.Update(ctx, &updatedAttendance)
		if err != nil {
			levelLogger.Error("Update failed", logger.String("error", err.Error()))
			t.Fatalf("Update failed: %v", err)
		}

		attendance, err := repo.GetByID(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			levelLogger.Error("GetByID after update failed", logger.String("error", err.Error()))
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if attendance.PresenceMark != false {
			levelLogger.Error("PresenceMark mismatch after update", logger.Bool("expected", false), logger.Bool("got", attendance.PresenceMark))
			t.Errorf("Expected PresenceMark false after update, got %v", attendance.PresenceMark)
		}
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{Field: "student_id", Operator: "=", Value: testAttendance.StudentID},
			},
			Limit: 10,
		}

		attendances, err := repo.List(ctx, filter)
		if err != nil {
			levelLogger.Error("List failed", logger.String("error", err.Error()))
			t.Fatalf("List failed: %v", err)
		}

		if len(attendances) == 0 {
			levelLogger.Error("Expected at least one attendance record in List")
			t.Error("Expected at least one attendance record in List")
		}
	})

	t.Run("GetByStudent", func(t *testing.T) {
		records, err := mgr.GetByStudent(ctx, testAttendance.StudentID)
		if err != nil {
			levelLogger.Error("GetByStudent failed", logger.String("error", err.Error()), logger.Int("studentID", testAttendance.StudentID))
		}
		assert.NoError(t, err)

		if len(records) == 0 {
			levelLogger.Error("Expected attendance records for student", logger.Int("studentID", testAttendance.StudentID))
		}
		assert.NotEmpty(t, records, "Expected attendance records for student")
	})

	t.Run("GetByLesson", func(t *testing.T) {
		records, err := mgr.GetByLesson(ctx, testAttendance.LessonID)
		if err != nil {
			levelLogger.Error("GetByLesson failed", logger.String("error", err.Error()), logger.Int("lessonID", testAttendance.LessonID))
		}
		assert.NoError(t, err)

		if len(records) == 0 {
			levelLogger.Error("Expected attendance records for lesson", logger.Int("lessonID", testAttendance.LessonID))
		}
		assert.NotEmpty(t, records, "Expected attendance records for lesson")
	})

	t.Run("GetByDateRange", func(t *testing.T) {
		startDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		endDate := time.Now().Format("2006-01-02")

		records, err := mgr.GetByDateRange(ctx, startDate, endDate)
		if err != nil {
			levelLogger.Error("GetByDateRange failed", logger.String("error", err.Error()), logger.String("startDate", startDate), logger.String("endDate", endDate))
		}
		assert.NoError(t, err)

		if len(records) == 0 {
			levelLogger.Error("Expected attendance records for date range", logger.String("startDate", startDate), logger.String("endDate", endDate))
		}
		assert.NotEmpty(t, records, "Expected attendance records for date range")
	})

	t.Run("GetStudentAttendanceStats", func(t *testing.T) {
		present, absent, err := mgr.GetStudentAttendanceStats(ctx, testAttendance.StudentID)
		if err != nil {
			levelLogger.Error("GetStudentAttendanceStats failed", logger.String("error", err.Error()), logger.Int("studentID", testAttendance.StudentID))
		}
		assert.NoError(t, err)

		// Проверяем что сумма присутствующих и отсутствующих >= 0
		if present < 0 || absent < 0 {
			levelLogger.Error("Invalid attendance stats", logger.Int("present", present), logger.Int("absent", absent))
			t.Error("Invalid attendance stats: negative counts")
		}
	})

	t.Run("BulkCreate", func(t *testing.T) {
		records := []*domain.StudentAttendance{
			{StudentID: 3, LessonID: 3, PresenceMark: true, AttendanceDate: time.Now().Format("2006-01-02")},
			{StudentID: 4, LessonID: 4, PresenceMark: false, AttendanceDate: time.Now().Format("2006-01-02")},
		}

		err := mgr.BulkCreate(ctx, records)
		if err != nil {
			levelLogger.Error("BulkCreate failed", logger.String("error", err.Error()))
		}
		assert.NoError(t, err)

		for _, record := range records {
			exists, err := repo.Exists(ctx, record.AttendanceNoteID)
			if err != nil {
				levelLogger.Error("Exists check failed after BulkCreate", logger.String("error", err.Error()), logger.Int("ID", record.AttendanceNoteID))
			}
			assert.NoError(t, err)

			if !exists {
				levelLogger.Error("Expected attendance to exist after BulkCreate", logger.Int("ID", record.AttendanceNoteID))
			}
			assert.True(t, exists, "Expected attendance to exist after BulkCreate")
		}
	})

	t.Run("CheckDuplicate", func(t *testing.T) {
		dup, err := mgr.CheckDuplicate(ctx, testAttendance.StudentID, testAttendance.LessonID)
		if err != nil {
			levelLogger.Error("CheckDuplicate failed", logger.String("error", err.Error()), logger.Int("studentID", testAttendance.StudentID), logger.Int("lessonID", testAttendance.LessonID))
		}
		assert.NoError(t, err)
		assert.True(t, dup, "Expected duplicate to be true for existing record")
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			levelLogger.Error("Delete failed", logger.String("error", err.Error()))
			t.Fatalf("Delete failed: %v", err)
		}

		exists, err := repo.Exists(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			levelLogger.Error("Exists after delete failed", logger.String("error", err.Error()))
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			levelLogger.Error("Expected attendance to be deleted")
			t.Error("Expected attendance to be deleted")
		}
	})

	if !t.Failed() {
		levelLogger.Info("All tests passed successfully")
		t.Log("All tests passed successfully")
	}
}

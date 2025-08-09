<<<<<<< HEAD
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

func TestStudentAttendanceRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentAttendanceRepository(sqlDB)

	// Тестовые данные
	testAttendance := &domain.StudentAttendance{
		StudentID:      2,
		LessonID:       2,
		PresenceMark:   true,
		AttendanceDate: time.Now().Format("2006-01-02"),
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testAttendance)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testAttendance.AttendanceNoteID == 0 {
			t.Error("Expected AttendanceNoteID to be set after Create")
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		attendance, err := repo.GetByID(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if attendance == nil {
			t.Fatal("Expected attendance to be found")
		}

		if attendance.StudentID != testAttendance.StudentID {
			t.Errorf("Expected StudentID %d, got %d", testAttendance.StudentID, attendance.StudentID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedAttendance := *testAttendance
		updatedAttendance.PresenceMark = false

		err = repoWithTx.Update(ctx, &updatedAttendance)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		attendance, err := repo.GetByID(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if attendance.PresenceMark != false {
			t.Errorf("Expected PresenceMark false after update, got %v", attendance.PresenceMark)
		}
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "student_id",
					Operator: "=",
					Value:    testAttendance.StudentID,
				},
			},
			Limit: 10,
		}

		attendances, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(attendances) == 0 {
			t.Error("Expected at least one attendance in List")
		}

		found := false
		for _, a := range attendances {
			if a.AttendanceNoteID == testAttendance.AttendanceNoteID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Created attendance not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "student_id",
					Operator: "=",
					Value:    testAttendance.StudentID,
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
		exists, err := repo.Exists(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}

		if !exists {
			t.Error("Expected attendance to exist")
		}

		// Проверка несуществующей записи
		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}

		if exists {
			t.Error("Expected non-existent attendance to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondAttendance := &domain.StudentAttendance{
			StudentID:      3,
			LessonID:       3,
			PresenceMark:   true,
			AttendanceDate: time.Now().Format("2006-01-02"),
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondAttendance)
		if err != nil {
			t.Fatalf("Create second attendance failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testAttendance.AttendanceNoteID, secondAttendance.AttendanceNoteID}
		attendances, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(attendances) != 2 {
			t.Errorf("Expected 2 attendances, got %d", len(attendances))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			t.Error("Expected attendance to be deleted")
		}
	})
}
=======
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

func TestStudentAttendanceRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentAttendanceRepository(sqlDB)

	// Тестовые данные
	testAttendance := &domain.StudentAttendance{
		StudentID:      2,
		LessonID:       2,
		PresenceMark:   true,
		AttendanceDate: time.Now().Format("2006-01-02"),
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testAttendance)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testAttendance.AttendanceNoteID == 0 {
			t.Error("Expected AttendanceNoteID to be set after Create")
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		attendance, err := repo.GetByID(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if attendance == nil {
			t.Fatal("Expected attendance to be found")
		}

		if attendance.StudentID != testAttendance.StudentID {
			t.Errorf("Expected StudentID %d, got %d", testAttendance.StudentID, attendance.StudentID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedAttendance := *testAttendance
		updatedAttendance.PresenceMark = false

		err = repoWithTx.Update(ctx, &updatedAttendance)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		attendance, err := repo.GetByID(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if attendance.PresenceMark != false {
			t.Errorf("Expected PresenceMark false after update, got %v", attendance.PresenceMark)
		}
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "student_id",
					Operator: "=",
					Value:    testAttendance.StudentID,
				},
			},
			Limit: 10,
		}

		attendances, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(attendances) == 0 {
			t.Error("Expected at least one attendance in List")
		}

		found := false
		for _, a := range attendances {
			if a.AttendanceNoteID == testAttendance.AttendanceNoteID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Created attendance not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "student_id",
					Operator: "=",
					Value:    testAttendance.StudentID,
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
		exists, err := repo.Exists(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}

		if !exists {
			t.Error("Expected attendance to exist")
		}

		// Проверка несуществующей записи
		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}

		if exists {
			t.Error("Expected non-existent attendance to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondAttendance := &domain.StudentAttendance{
			StudentID:      3,
			LessonID:       3,
			PresenceMark:   true,
			AttendanceDate: time.Now().Format("2006-01-02"),
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondAttendance)
		if err != nil {
			t.Fatalf("Create second attendance failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testAttendance.AttendanceNoteID, secondAttendance.AttendanceNoteID}
		attendances, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(attendances) != 2 {
			t.Errorf("Expected 2 attendances, got %d", len(attendances))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testAttendance.AttendanceNoteID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			t.Error("Expected attendance to be deleted")
		}
	})
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

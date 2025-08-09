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

func TestSubjectRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewSubjectRepository(sqlDB)

	testSubject := &domain.Subject{
		SubjectName: "Музыка",
		SubjectType: "Теория",
		ShortDesc:   "Краткое описание предмета",
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testSubject)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testSubject.SubjectID == 0 {
			t.Error("Expected SubjectID to be set after Create")
		} else {
			t.Logf("Created SubjectID: %d", testSubject.SubjectID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		subj, err := repo.GetByID(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if subj == nil {
			t.Fatal("Expected Subject to be found")
		}
		if subj.SubjectName != testSubject.SubjectName {
			t.Errorf("Expected SubjectName %q, got %q", testSubject.SubjectName, subj.SubjectName)
		}
		if subj.SubjectType != testSubject.SubjectType {
			t.Errorf("Expected SubjectType %q, got %q", testSubject.SubjectType, subj.SubjectType)
		}
		if subj.ShortDesc != testSubject.ShortDesc {
			t.Errorf("Expected ShortDesc %q, got %q", testSubject.ShortDesc, subj.ShortDesc)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedSubject := *testSubject
		updatedSubject.SubjectName = "Музыкальная теория"
		updatedSubject.SubjectType = "Практика"
		updatedSubject.ShortDesc = "Обновленное описание"

		err = repoWithTx.Update(ctx, &updatedSubject)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		subj, err := repo.GetByID(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if subj.SubjectName != updatedSubject.SubjectName {
			t.Errorf("Expected SubjectName %q after update, got %q", updatedSubject.SubjectName, subj.SubjectName)
		}
		if subj.SubjectType != updatedSubject.SubjectType {
			t.Errorf("Expected SubjectType %q after update, got %q", updatedSubject.SubjectType, subj.SubjectType)
		}
		if subj.ShortDesc != updatedSubject.ShortDesc {
			t.Errorf("Expected ShortDesc %q after update, got %q", updatedSubject.ShortDesc, subj.ShortDesc)
		}

		*testSubject = updatedSubject
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "subject_name",
					Operator: "=",
					Value:    testSubject.SubjectName,
				},
			},
			Limit: 10,
		}

		subjects, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(subjects) == 0 {
			t.Error("Expected at least one Subject in List")
		} else {
			t.Logf("List returned %d items", len(subjects))
		}

		found := false
		for _, s := range subjects {
			if s.SubjectID == testSubject.SubjectID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Subject not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "subject_name",
					Operator: "=",
					Value:    testSubject.SubjectName,
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
		exists, err := repo.Exists(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Subject to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Subject to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondSubject := &domain.Subject{
			SubjectName: "История музыки",
			SubjectType: "Лекция",
			ShortDesc:   "Описание второго предмета",
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondSubject)
		if err != nil {
			t.Fatalf("Create second Subject failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testSubject.SubjectID, secondSubject.SubjectID}
		subjects, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(subjects) != 2 {
			t.Errorf("Expected 2 Subjects, got %d", len(subjects))
		}
	})

	t.Run("GetPopularSubjects", func(t *testing.T) {
		// Для корректного теста нужна реальная связанная таблица programm_distributions с данными
		limit := 5
		popularSubjects, err := repo.GetPopularSubjects(ctx, limit)
		if err != nil {
			t.Fatalf("GetPopularSubjects failed: %v", err)
		}

		if len(popularSubjects) > limit {
			t.Errorf("Expected at most %d subjects, got %d", limit, len(popularSubjects))
		}

		t.Logf("GetPopularSubjects returned %d subjects", len(popularSubjects))
	})

	t.Run("GetSubjectsWithPrograms", func(t *testing.T) {
		// Для корректного теста нужен валидный musprogramm_id в programm_distributions
		programID := 1
		subjectsWithProgram, err := repo.GetSubjectsWithPrograms(ctx, programID)
		if err != nil {
			t.Fatalf("GetSubjectsWithPrograms failed: %v", err)
		}

		t.Logf("GetSubjectsWithPrograms returned %d subjects for programID %d", len(subjectsWithProgram), programID)
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Subject to be deleted")
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

func TestSubjectRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewSubjectRepository(sqlDB)

	testSubject := &domain.Subject{
		SubjectName: "Музыка",
		SubjectType: "Теория",
		ShortDesc:   "Краткое описание предмета",
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testSubject)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testSubject.SubjectID == 0 {
			t.Error("Expected SubjectID to be set after Create")
		} else {
			t.Logf("Created SubjectID: %d", testSubject.SubjectID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		subj, err := repo.GetByID(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if subj == nil {
			t.Fatal("Expected Subject to be found")
		}
		if subj.SubjectName != testSubject.SubjectName {
			t.Errorf("Expected SubjectName %q, got %q", testSubject.SubjectName, subj.SubjectName)
		}
		if subj.SubjectType != testSubject.SubjectType {
			t.Errorf("Expected SubjectType %q, got %q", testSubject.SubjectType, subj.SubjectType)
		}
		if subj.ShortDesc != testSubject.ShortDesc {
			t.Errorf("Expected ShortDesc %q, got %q", testSubject.ShortDesc, subj.ShortDesc)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedSubject := *testSubject
		updatedSubject.SubjectName = "Музыкальная теория"
		updatedSubject.SubjectType = "Практика"
		updatedSubject.ShortDesc = "Обновленное описание"

		err = repoWithTx.Update(ctx, &updatedSubject)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		subj, err := repo.GetByID(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if subj.SubjectName != updatedSubject.SubjectName {
			t.Errorf("Expected SubjectName %q after update, got %q", updatedSubject.SubjectName, subj.SubjectName)
		}
		if subj.SubjectType != updatedSubject.SubjectType {
			t.Errorf("Expected SubjectType %q after update, got %q", updatedSubject.SubjectType, subj.SubjectType)
		}
		if subj.ShortDesc != updatedSubject.ShortDesc {
			t.Errorf("Expected ShortDesc %q after update, got %q", updatedSubject.ShortDesc, subj.ShortDesc)
		}

		*testSubject = updatedSubject
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "subject_name",
					Operator: "=",
					Value:    testSubject.SubjectName,
				},
			},
			Limit: 10,
		}

		subjects, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(subjects) == 0 {
			t.Error("Expected at least one Subject in List")
		} else {
			t.Logf("List returned %d items", len(subjects))
		}

		found := false
		for _, s := range subjects {
			if s.SubjectID == testSubject.SubjectID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Subject not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "subject_name",
					Operator: "=",
					Value:    testSubject.SubjectName,
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
		exists, err := repo.Exists(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Subject to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Subject to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondSubject := &domain.Subject{
			SubjectName: "История музыки",
			SubjectType: "Лекция",
			ShortDesc:   "Описание второго предмета",
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondSubject)
		if err != nil {
			t.Fatalf("Create second Subject failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testSubject.SubjectID, secondSubject.SubjectID}
		subjects, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(subjects) != 2 {
			t.Errorf("Expected 2 Subjects, got %d", len(subjects))
		}
	})

	t.Run("GetPopularSubjects", func(t *testing.T) {
		// Для корректного теста нужна реальная связанная таблица programm_distributions с данными
		limit := 5
		popularSubjects, err := repo.GetPopularSubjects(ctx, limit)
		if err != nil {
			t.Fatalf("GetPopularSubjects failed: %v", err)
		}

		if len(popularSubjects) > limit {
			t.Errorf("Expected at most %d subjects, got %d", limit, len(popularSubjects))
		}

		t.Logf("GetPopularSubjects returned %d subjects", len(popularSubjects))
	})

	t.Run("GetSubjectsWithPrograms", func(t *testing.T) {
		// Для корректного теста нужен валидный musprogramm_id в programm_distributions
		programID := 1
		subjectsWithProgram, err := repo.GetSubjectsWithPrograms(ctx, programID)
		if err != nil {
			t.Fatalf("GetSubjectsWithPrograms failed: %v", err)
		}

		t.Logf("GetSubjectsWithPrograms returned %d subjects for programID %d", len(subjectsWithProgram), programID)
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testSubject.SubjectID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Subject to be deleted")
		}
	})
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

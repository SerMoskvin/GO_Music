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

func TestStudyGroupRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudyGroupRepository(sqlDB)

	testGroup := &domain.StudyGroup{
		MusProgrammID:    1, // валидный ID музыкальной программы
		GroupName:        "Группа A",
		StudyYear:        2023,
		NumberOfStudents: 25,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testGroup)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testGroup.GroupID == 0 {
			t.Error("Expected GroupID to be set after Create")
		} else {
			t.Logf("Created GroupID: %d", testGroup.GroupID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		group, err := repo.GetByID(ctx, testGroup.GroupID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if group == nil {
			t.Fatal("Expected StudyGroup to be found")
		}
		if group.GroupName != testGroup.GroupName {
			t.Errorf("Expected GroupName %q, got %q", testGroup.GroupName, group.GroupName)
		}
		if group.MusProgrammID != testGroup.MusProgrammID {
			t.Errorf("Expected MusProgrammID %d, got %d", testGroup.MusProgrammID, group.MusProgrammID)
		}
		if group.StudyYear != testGroup.StudyYear {
			t.Errorf("Expected StudyYear %d, got %d", testGroup.StudyYear, group.StudyYear)
		}
		if group.NumberOfStudents != testGroup.NumberOfStudents {
			t.Errorf("Expected NumberOfStudents %d, got %d", testGroup.NumberOfStudents, group.NumberOfStudents)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedGroup := *testGroup
		updatedGroup.GroupName = "Группа Б"
		updatedGroup.StudyYear = 2024
		updatedGroup.NumberOfStudents = 30

		err = repoWithTx.Update(ctx, &updatedGroup)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		group, err := repo.GetByID(ctx, testGroup.GroupID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if group.GroupName != "Группа Б" {
			t.Errorf("Expected GroupName 'Группа Б' after update, got %q", group.GroupName)
		}
		if group.StudyYear != 2024 {
			t.Errorf("Expected StudyYear 2024 after update, got %d", group.StudyYear)
		}
		if group.NumberOfStudents != 30 {
			t.Errorf("Expected NumberOfStudents 30 after update, got %d", group.NumberOfStudents)
		}

		*testGroup = updatedGroup
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "group_name",
					Operator: "=",
					Value:    testGroup.GroupName,
				},
			},
			Limit: 10,
		}

		groups, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(groups) == 0 {
			t.Error("Expected at least one StudyGroup in List")
		} else {
			t.Logf("List returned %d items", len(groups))
		}

		found := false
		for _, g := range groups {
			if g.GroupID == testGroup.GroupID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created StudyGroup not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "group_name",
					Operator: "=",
					Value:    testGroup.GroupName,
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
		exists, err := repo.Exists(ctx, testGroup.GroupID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected StudyGroup to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent StudyGroup to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondGroup := &domain.StudyGroup{
			MusProgrammID:    1,
			GroupName:        "Группа В",
			StudyYear:        2023,
			NumberOfStudents: 20,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondGroup)
		if err != nil {
			t.Fatalf("Create second StudyGroup failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testGroup.GroupID, secondGroup.GroupID}
		groups, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(groups) != 2 {
			t.Errorf("Expected 2 StudyGroups, got %d", len(groups))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testGroup.GroupID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testGroup.GroupID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected StudyGroup to be deleted")
		}
	})
}

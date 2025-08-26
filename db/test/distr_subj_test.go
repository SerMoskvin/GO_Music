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

func TestSubjectDistributionRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewSubjectDistributionRepository(sqlDB)

	testSD := &domain.SubjectDistribution{
		EmployeeID: 6, // укажите валидные ID существующих записей
		SubjectID:  2,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testSD)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testSD.SubjectDistrID == 0 {
			t.Error("Expected SubjectDistrID to be set after Create")
		} else {
			t.Logf("Created SubjectDistrID: %d", testSD.SubjectDistrID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		sd, err := repo.GetByID(ctx, testSD.SubjectDistrID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if sd == nil {
			t.Fatal("Expected SubjectDistribution to be found")
		}
		if sd.EmployeeID != testSD.EmployeeID {
			t.Errorf("Expected EmployeeID %d, got %d", testSD.EmployeeID, sd.EmployeeID)
		}
		if sd.SubjectID != testSD.SubjectID {
			t.Errorf("Expected SubjectID %d, got %d", testSD.SubjectID, sd.SubjectID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedSD := *testSD
		updatedSD.EmployeeID = 7
		updatedSD.SubjectID = 4

		err = repoWithTx.Update(ctx, &updatedSD)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		sd, err := repo.GetByID(ctx, testSD.SubjectDistrID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if sd.EmployeeID != 7 {
			t.Errorf("Expected EmployeeID 7 after update, got %d", sd.EmployeeID)
		}
		if sd.SubjectID != 4 {
			t.Errorf("Expected SubjectID 4 after update, got %d", sd.SubjectID)
		}

		*testSD = updatedSD
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "employee_id",
					Operator: "=",
					Value:    testSD.EmployeeID,
				},
			},
			Limit: 10,
		}

		sds, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(sds) == 0 {
			t.Error("Expected at least one SubjectDistribution in List")
		} else {
			t.Logf("List returned %d items", len(sds))
		}

		found := false
		for _, sd := range sds {
			if sd.SubjectDistrID == testSD.SubjectDistrID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created SubjectDistribution not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "employee_id",
					Operator: "=",
					Value:    testSD.EmployeeID,
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
		exists, err := repo.Exists(ctx, testSD.SubjectDistrID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected SubjectDistribution to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent SubjectDistribution to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondSD := &domain.SubjectDistribution{
			EmployeeID: 8,
			SubjectID:  2,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondSD)
		if err != nil {
			t.Fatalf("Create second SubjectDistribution failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testSD.SubjectDistrID, secondSD.SubjectDistrID}
		sds, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(sds) != 2 {
			t.Errorf("Expected 2 SubjectDistributions, got %d", len(sds))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testSD.SubjectDistrID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testSD.SubjectDistrID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected SubjectDistribution to be deleted")
		}
	})
}

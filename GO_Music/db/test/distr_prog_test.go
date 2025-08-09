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

func TestProgrammDistributionRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewProgrammDistributionRepository(sqlDB)

	testPD := &domain.ProgrammDistribution{
		MusprogrammID: 1,
		SubjectID:     2,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testPD)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testPD.ProgrammDistrID == 0 {
			t.Error("Expected ProgrammDistrID to be set after Create")
		} else {
			t.Logf("Created ProgrammDistrID: %d", testPD.ProgrammDistrID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		pd, err := repo.GetByID(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if pd == nil {
			t.Fatal("Expected ProgrammDistribution to be found")
		}
		if pd.MusprogrammID != testPD.MusprogrammID {
			t.Errorf("Expected MusprogrammID %d, got %d", testPD.MusprogrammID, pd.MusprogrammID)
		}
		if pd.SubjectID != testPD.SubjectID {
			t.Errorf("Expected SubjectID %d, got %d", testPD.SubjectID, pd.SubjectID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedPD := *testPD
		updatedPD.MusprogrammID = 3
		updatedPD.SubjectID = 4

		err = repoWithTx.Update(ctx, &updatedPD)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		pd, err := repo.GetByID(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if pd.MusprogrammID != 3 {
			t.Errorf("Expected MusprogrammID 3 after update, got %d", pd.MusprogrammID)
		}
		if pd.SubjectID != 4 {
			t.Errorf("Expected SubjectID 4 after update, got %d", pd.SubjectID)
		}

		*testPD = updatedPD
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "musprogramm_id",
					Operator: "=",
					Value:    testPD.MusprogrammID,
				},
			},
			Limit: 10,
		}

		pds, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(pds) == 0 {
			t.Error("Expected at least one ProgrammDistribution in List")
		} else {
			t.Logf("List returned %d items", len(pds))
		}

		found := false
		for _, pd := range pds {
			if pd.ProgrammDistrID == testPD.ProgrammDistrID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created ProgrammDistribution not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "musprogramm_id",
					Operator: "=",
					Value:    testPD.MusprogrammID,
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
		exists, err := repo.Exists(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected ProgrammDistribution to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent ProgrammDistribution to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondPD := &domain.ProgrammDistribution{
			MusprogrammID: 5,
			SubjectID:     5,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondPD)
		if err != nil {
			t.Fatalf("Create second ProgrammDistribution failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testPD.ProgrammDistrID, secondPD.ProgrammDistrID}
		pds, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(pds) != 2 {
			t.Errorf("Expected 2 ProgrammDistributions, got %d", len(pds))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected ProgrammDistribution to be deleted")
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

func TestProgrammDistributionRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewProgrammDistributionRepository(sqlDB)

	testPD := &domain.ProgrammDistribution{
		MusprogrammID: 1,
		SubjectID:     2,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testPD)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testPD.ProgrammDistrID == 0 {
			t.Error("Expected ProgrammDistrID to be set after Create")
		} else {
			t.Logf("Created ProgrammDistrID: %d", testPD.ProgrammDistrID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		pd, err := repo.GetByID(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if pd == nil {
			t.Fatal("Expected ProgrammDistribution to be found")
		}
		if pd.MusprogrammID != testPD.MusprogrammID {
			t.Errorf("Expected MusprogrammID %d, got %d", testPD.MusprogrammID, pd.MusprogrammID)
		}
		if pd.SubjectID != testPD.SubjectID {
			t.Errorf("Expected SubjectID %d, got %d", testPD.SubjectID, pd.SubjectID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedPD := *testPD
		updatedPD.MusprogrammID = 3
		updatedPD.SubjectID = 4

		err = repoWithTx.Update(ctx, &updatedPD)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		pd, err := repo.GetByID(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if pd.MusprogrammID != 3 {
			t.Errorf("Expected MusprogrammID 3 after update, got %d", pd.MusprogrammID)
		}
		if pd.SubjectID != 4 {
			t.Errorf("Expected SubjectID 4 after update, got %d", pd.SubjectID)
		}

		*testPD = updatedPD
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "musprogramm_id",
					Operator: "=",
					Value:    testPD.MusprogrammID,
				},
			},
			Limit: 10,
		}

		pds, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(pds) == 0 {
			t.Error("Expected at least one ProgrammDistribution in List")
		} else {
			t.Logf("List returned %d items", len(pds))
		}

		found := false
		for _, pd := range pds {
			if pd.ProgrammDistrID == testPD.ProgrammDistrID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created ProgrammDistribution not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "musprogramm_id",
					Operator: "=",
					Value:    testPD.MusprogrammID,
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
		exists, err := repo.Exists(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected ProgrammDistribution to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent ProgrammDistribution to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondPD := &domain.ProgrammDistribution{
			MusprogrammID: 5,
			SubjectID:     5,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondPD)
		if err != nil {
			t.Fatalf("Create second ProgrammDistribution failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testPD.ProgrammDistrID, secondPD.ProgrammDistrID}
		pds, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(pds) != 2 {
			t.Errorf("Expected 2 ProgrammDistributions, got %d", len(pds))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testPD.ProgrammDistrID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected ProgrammDistribution to be deleted")
		}
	})
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

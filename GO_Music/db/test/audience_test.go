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

func TestAudienceRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewAudienceRepository(sqlDB)

	testAudience := &domain.Audience{
		Name:        "Test Audience",
		AudinType:   "Lecture Hall",
		AudinNumber: "101",
		Capacity:    50,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testAudience)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testAudience.AudienceID == 0 {
			t.Error("Expected AudienceID to be set after Create")
		} else {
			t.Logf("Created AudienceID: %d", testAudience.AudienceID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		audience, err := repo.GetByID(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if audience == nil {
			t.Fatal("Expected audience to be found")
		}

		if audience.Name != testAudience.Name {
			t.Errorf("Expected Name %q, got %q", testAudience.Name, audience.Name)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedAudience := *testAudience
		updatedAudience.Name = "Updated Audience"
		updatedAudience.Capacity = 100

		err = repoWithTx.Update(ctx, &updatedAudience)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		audience, err := repo.GetByID(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if audience.Name != "Updated Audience" {
			t.Errorf("Expected Name 'Updated Audience' after update, got %q", audience.Name)
		}
		if audience.Capacity != 100 {
			t.Errorf("Expected Capacity 100 after update, got %d", audience.Capacity)
		}

		*testAudience = updatedAudience
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "name",
					Operator: "=",
					Value:    testAudience.Name,
				},
			},
			Limit: 10,
		}

		audiences, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(audiences) == 0 {
			t.Error("Expected at least one audience in List")
		} else {
			t.Logf("List returned %d audiences", len(audiences))
		}

		found := false
		for _, a := range audiences {
			if a.AudienceID == testAudience.AudienceID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Created audience not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "name",
					Operator: "=",
					Value:    testAudience.Name,
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
		exists, err := repo.Exists(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}

		if !exists {
			t.Error("Expected audience to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}

		if exists {
			t.Error("Expected non-existent audience to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondAudience := &domain.Audience{
			Name:        "Second Audience",
			AudinType:   "Seminar Room",
			AudinNumber: "202",
			Capacity:    30,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondAudience)
		if err != nil {
			t.Fatalf("Create second audience failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testAudience.AudienceID, secondAudience.AudienceID}
		audiences, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(audiences) != 2 {
			t.Errorf("Expected 2 audiences, got %d", len(audiences))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			t.Error("Expected audience to be deleted")
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

func TestAudienceRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewAudienceRepository(sqlDB)

	testAudience := &domain.Audience{
		Name:        "Test Audience",
		AudinType:   "Lecture Hall",
		AudinNumber: "101",
		Capacity:    50,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testAudience)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testAudience.AudienceID == 0 {
			t.Error("Expected AudienceID to be set after Create")
		} else {
			t.Logf("Created AudienceID: %d", testAudience.AudienceID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		audience, err := repo.GetByID(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}

		if audience == nil {
			t.Fatal("Expected audience to be found")
		}

		if audience.Name != testAudience.Name {
			t.Errorf("Expected Name %q, got %q", testAudience.Name, audience.Name)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedAudience := *testAudience
		updatedAudience.Name = "Updated Audience"
		updatedAudience.Capacity = 100

		err = repoWithTx.Update(ctx, &updatedAudience)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		audience, err := repo.GetByID(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}

		if audience.Name != "Updated Audience" {
			t.Errorf("Expected Name 'Updated Audience' after update, got %q", audience.Name)
		}
		if audience.Capacity != 100 {
			t.Errorf("Expected Capacity 100 after update, got %d", audience.Capacity)
		}

		*testAudience = updatedAudience
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "name",
					Operator: "=",
					Value:    testAudience.Name,
				},
			},
			Limit: 10,
		}

		audiences, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(audiences) == 0 {
			t.Error("Expected at least one audience in List")
		} else {
			t.Logf("List returned %d audiences", len(audiences))
		}

		found := false
		for _, a := range audiences {
			if a.AudienceID == testAudience.AudienceID {
				found = true
				break
			}
		}

		if !found {
			t.Error("Created audience not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "name",
					Operator: "=",
					Value:    testAudience.Name,
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
		exists, err := repo.Exists(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}

		if !exists {
			t.Error("Expected audience to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}

		if exists {
			t.Error("Expected non-existent audience to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondAudience := &domain.Audience{
			Name:        "Second Audience",
			AudinType:   "Seminar Room",
			AudinNumber: "202",
			Capacity:    30,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondAudience)
		if err != nil {
			t.Fatalf("Create second audience failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testAudience.AudienceID, secondAudience.AudienceID}
		audiences, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(audiences) != 2 {
			t.Errorf("Expected 2 audiences, got %d", len(audiences))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testAudience.AudienceID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}

		if exists {
			t.Error("Expected audience to be deleted")
		}
	})
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

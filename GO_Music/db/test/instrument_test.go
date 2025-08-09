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

func TestInstrumentRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewInstrumentRepository(sqlDB)

	testInstrument := &domain.Instrument{
		AudienceID: 1, // Укажите валидный ID аудитории
		Name:       "Фортепиано",
		InstrType:  "Клавишный",
		Condition:  "Хорошее",
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testInstrument)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testInstrument.InstrumentID == 0 {
			t.Error("Expected InstrumentID to be set after Create")
		} else {
			t.Logf("Created InstrumentID: %d", testInstrument.InstrumentID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		inst, err := repo.GetByID(ctx, testInstrument.InstrumentID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if inst == nil {
			t.Fatal("Expected Instrument to be found")
		}
		if inst.Name != testInstrument.Name {
			t.Errorf("Expected Name %q, got %q", testInstrument.Name, inst.Name)
		}
		if inst.AudienceID != testInstrument.AudienceID {
			t.Errorf("Expected AudienceID %d, got %d", testInstrument.AudienceID, inst.AudienceID)
		}
		if inst.InstrType != testInstrument.InstrType {
			t.Errorf("Expected InstrType %q, got %q", testInstrument.InstrType, inst.InstrType)
		}
		if inst.Condition != testInstrument.Condition {
			t.Errorf("Expected Condition %q, got %q", testInstrument.Condition, inst.Condition)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedInst := *testInstrument
		updatedInst.Name = "Скрипка"
		updatedInst.InstrType = "Струнный"
		updatedInst.Condition = "Отличное"

		err = repoWithTx.Update(ctx, &updatedInst)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		inst, err := repo.GetByID(ctx, testInstrument.InstrumentID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if inst.Name != "Скрипка" {
			t.Errorf("Expected Name 'Скрипка' after update, got %q", inst.Name)
		}
		if inst.InstrType != "Струнный" {
			t.Errorf("Expected InstrType 'Струнный' after update, got %q", inst.InstrType)
		}
		if inst.Condition != "Отличное" {
			t.Errorf("Expected Condition 'Отличное' after update, got %q", inst.Condition)
		}

		*testInstrument = updatedInst
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "name",
					Operator: "=",
					Value:    testInstrument.Name,
				},
			},
			Limit: 10,
		}

		instruments, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(instruments) == 0 {
			t.Error("Expected at least one Instrument in List")
		} else {
			t.Logf("List returned %d items", len(instruments))
		}

		found := false
		for _, i := range instruments {
			if i.InstrumentID == testInstrument.InstrumentID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Instrument not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "name",
					Operator: "=",
					Value:    testInstrument.Name,
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
		exists, err := repo.Exists(ctx, testInstrument.InstrumentID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Instrument to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Instrument to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondInst := &domain.Instrument{
			AudienceID: 1,
			Name:       "Гитара",
			InstrType:  "Струнный",
			Condition:  "Хорошее",
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondInst)
		if err != nil {
			t.Fatalf("Create second Instrument failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testInstrument.InstrumentID, secondInst.InstrumentID}
		instruments, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(instruments) != 2 {
			t.Errorf("Expected 2 Instruments, got %d", len(instruments))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testInstrument.InstrumentID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testInstrument.InstrumentID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Instrument to be deleted")
		}
	})
}

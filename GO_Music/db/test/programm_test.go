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

func TestProgrammRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewProgrammRepository(sqlDB)

	// Пример nullable строк
	instrument := "Фортепиано"
	description := "Описание программы"
	point := &instrument

	testProgramm := &domain.Programm{
		ProgrammName:           "Музыкальная программа 1",
		ProgrammType:           "Тип программы",
		Duration:               120,
		Instrument:             point,
		Description:            &description,
		StudyLoad:              60,
		FinalCertificationForm: "Экзамен",
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		// Логируем объект перед созданием
		t.Logf("Creating Programm: %+v", testProgramm)

		err = repoWithTx.Create(ctx, testProgramm)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testProgramm.MusprogrammID == 0 {
			t.Error("Expected MusprogrammID to be set after Create")
		} else {
			t.Logf("Created MusprogrammID: %d", testProgramm.MusprogrammID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
		var instrumentFromDB, descriptionFromDB *string
		err = sqlDB.QueryRowContext(ctx, "SELECT instrument, description FROM programm WHERE musprogramm_id = $1", testProgramm.MusprogrammID).Scan(&instrumentFromDB, &descriptionFromDB)
		if err != nil {
			t.Fatalf("failed to query instrument and description from db: %v", err)
		}
		t.Logf("Instrument in DB after create: %v", instrumentFromDB)
		t.Logf("Description in DB after create: %v", descriptionFromDB)
	})

	t.Run("GetByID", func(t *testing.T) {
		p, err := repo.GetByID(ctx, testProgramm.MusprogrammID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if p == nil {
			t.Fatal("Expected Programm to be found")
		}

		// Логируем объект, полученный из БД
		t.Logf("Programm from DB: %+v", p)

		if p.ProgrammName != testProgramm.ProgrammName {
			t.Errorf("Expected ProgrammName %q, got %q", testProgramm.ProgrammName, p.ProgrammName)
		}
		if p.ProgrammType != testProgramm.ProgrammType {
			t.Errorf("Expected ProgrammType %q, got %q", testProgramm.ProgrammType, p.ProgrammType)
		}
		if p.Duration != testProgramm.Duration {
			t.Errorf("Expected Duration %d, got %d", testProgramm.Duration, p.Duration)
		}
		if (p.Instrument == nil) != (testProgramm.Instrument == nil) ||
			(p.Instrument != nil && *p.Instrument != *testProgramm.Instrument) {
			t.Errorf("Expected Instrument %v, got %v", testProgramm.Instrument, p.Instrument)
		}
		if (p.Description == nil) != (testProgramm.Description == nil) ||
			(p.Description != nil && *p.Description != *testProgramm.Description) {
			t.Errorf("Expected Description %v, got %v", testProgramm.Description, p.Description)
		}
		if p.StudyLoad != testProgramm.StudyLoad {
			t.Errorf("Expected StudyLoad %d, got %d", testProgramm.StudyLoad, p.StudyLoad)
		}
		if p.FinalCertificationForm != testProgramm.FinalCertificationForm {
			t.Errorf("Expected FinalCertificationForm %q, got %q", testProgramm.FinalCertificationForm, p.FinalCertificationForm)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		newInstrument := "Скрипка"
		newDescription := "Новое описание"
		updatedProgramm := *testProgramm
		updatedProgramm.ProgrammName = "Обновленная программа"
		updatedProgramm.ProgrammType = "Новый тип"
		updatedProgramm.Duration = 150
		updatedProgramm.Instrument = &newInstrument
		updatedProgramm.Description = &newDescription
		updatedProgramm.StudyLoad = 75
		updatedProgramm.FinalCertificationForm = "Зачет"

		t.Logf("Updating Programm to: %+v", updatedProgramm)

		err = repoWithTx.Update(ctx, &updatedProgramm)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		p, err := repo.GetByID(ctx, testProgramm.MusprogrammID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		t.Logf("Programm from DB after update: %+v", p)

		// сравнения как было
		if p.ProgrammName != updatedProgramm.ProgrammName {
			t.Errorf("Expected ProgrammName %q after update, got %q", updatedProgramm.ProgrammName, p.ProgrammName)
		}
		if p.ProgrammType != updatedProgramm.ProgrammType {
			t.Errorf("Expected ProgrammType %q after update, got %q", updatedProgramm.ProgrammType, p.ProgrammType)
		}
		if p.Duration != updatedProgramm.Duration {
			t.Errorf("Expected Duration %d after update, got %d", updatedProgramm.Duration, p.Duration)
		}
		if (p.Instrument == nil) != (updatedProgramm.Instrument == nil) ||
			(p.Instrument != nil && *p.Instrument != *updatedProgramm.Instrument) {
			t.Errorf("Expected Instrument %v after update, got %v", updatedProgramm.Instrument, p.Instrument)
		}
		if (p.Description == nil) != (updatedProgramm.Description == nil) ||
			(p.Description != nil && *p.Description != *updatedProgramm.Description) {
			t.Errorf("Expected Description %v after update, got %v", updatedProgramm.Description, p.Description)
		}
		if p.StudyLoad != updatedProgramm.StudyLoad {
			t.Errorf("Expected StudyLoad %d after update, got %d", updatedProgramm.StudyLoad, p.StudyLoad)
		}
		if p.FinalCertificationForm != updatedProgramm.FinalCertificationForm {
			t.Errorf("Expected FinalCertificationForm %q after update, got %q", updatedProgramm.FinalCertificationForm, p.FinalCertificationForm)
		}

		*testProgramm = updatedProgramm
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "programm_name",
					Operator: "=",
					Value:    testProgramm.ProgrammName,
				},
			},
			Limit: 10,
		}

		programms, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(programms) == 0 {
			t.Error("Expected at least one Programm in List")
		} else {
			t.Logf("List returned %d items", len(programms))
		}

		found := false
		for _, p := range programms {
			if p.MusprogrammID == testProgramm.MusprogrammID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Programm not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "programm_name",
					Operator: "=",
					Value:    testProgramm.ProgrammName,
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
		exists, err := repo.Exists(ctx, testProgramm.MusprogrammID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Programm to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Programm to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondProgramm := &domain.Programm{
			ProgrammName:           "Вторая программа",
			ProgrammType:           "Тип 2",
			Duration:               90,
			Instrument:             nil,
			Description:            nil,
			StudyLoad:              45,
			FinalCertificationForm: "Зачет",
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondProgramm)
		if err != nil {
			t.Fatalf("Create second Programm failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testProgramm.MusprogrammID, secondProgramm.MusprogrammID}
		programms, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(programms) != 2 {
			t.Errorf("Expected 2 Programms, got %d", len(programms))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testProgramm.MusprogrammID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testProgramm.MusprogrammID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Programm to be deleted")
		}
	})
}

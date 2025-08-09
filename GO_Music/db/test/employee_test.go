<<<<<<< HEAD
package db_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	_ "github.com/lib/pq"
)

func TestEmployeeRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewEmployeeRepository(sqlDB)
	FatherName := "Loxovich"

	testEmp := &domain.Employee{
		UserID:         nil,
		Surname:        "Иванов",
		Name:           "Иван",
		FatherName:     &FatherName,
		Birthday:       time.Date(1990, 5, 20, 0, 0, 0, 0, time.UTC),
		PhoneNumber:    "79991112233",
		Job:            "Программист",
		WorkExperience: 5,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testEmp)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testEmp.EmployeeID == 0 {
			t.Error("Expected EmployeeID to be set after Create")
		} else {
			t.Logf("Created EmployeeID: %d", testEmp.EmployeeID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
		t.Logf("testEmp.EmployeeID = %d", testEmp.EmployeeID)
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Logf("Trying to get employee with ID: %d", testEmp.EmployeeID)
		emp, err := repo.GetByID(ctx, testEmp.EmployeeID)
		t.Logf("Successfully retrieved employee: %+v", emp)

		// Проверки полей
		if emp.Surname != testEmp.Surname {
			t.Errorf("Expected Surname %q, got %q", testEmp.Surname, emp.Surname)
		}
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if emp == nil {
			t.Fatal("Expected Employee to be found")
		}
		if emp.Surname != testEmp.Surname {
			t.Errorf("Expected Surname %q, got %q", testEmp.Surname, emp.Surname)
		}
		if emp.Name != testEmp.Name {
			t.Errorf("Expected Name %q, got %q", testEmp.Name, emp.Name)
		}
		if !emp.Birthday.Equal(testEmp.Birthday) {
			t.Errorf("Expected Birthday %v, got %v", testEmp.Birthday, emp.Birthday)
		}
		if emp.PhoneNumber != testEmp.PhoneNumber {
			t.Errorf("Expected PhoneNumber %q, got %q", testEmp.PhoneNumber, emp.PhoneNumber)
		}
		if emp.Job != testEmp.Job {
			t.Errorf("Expected Job %q, got %q", testEmp.Job, emp.Job)
		}
		if emp.WorkExperience != testEmp.WorkExperience {
			t.Errorf("Expected WorkExperience %d, got %d", testEmp.WorkExperience, emp.WorkExperience)
		}
		fmt.Printf("Original: %+v", testEmp)
		fmt.Printf("FromDB: %+v", emp)
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedEmp := *testEmp
		updatedEmp.Surname = "Петров"
		updatedEmp.Name = "Пётр"
		updatedEmp.PhoneNumber = "78882223344"
		updatedEmp.Job = "Тестировщик"
		updatedEmp.WorkExperience = 7

		err = repoWithTx.Update(ctx, &updatedEmp)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		emp, err := repo.GetByID(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if emp.Surname != "Петров" {
			t.Errorf("Expected Surname 'Петров' after update, got %q", emp.Surname)
		}
		if emp.Name != "Пётр" {
			t.Errorf("Expected Name 'Пётр' after update, got %q", emp.Name)
		}
		if emp.PhoneNumber != "78882223344" {
			t.Errorf("Expected PhoneNumber '78882223344' after update, got %q", emp.PhoneNumber)
		}
		if emp.Job != "Тестировщик" {
			t.Errorf("Expected Job 'Тестировщик' after update, got %q", emp.Job)
		}
		if emp.WorkExperience != 7 {
			t.Errorf("Expected WorkExperience 7 after update, got %d", emp.WorkExperience)
		}

		*testEmp = updatedEmp
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testEmp.Surname,
				},
			},
			Limit: 10,
		}

		emps, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(emps) == 0 {
			t.Error("Expected at least one Employee in List")
		} else {
			t.Logf("List returned %d items", len(emps))
		}

		found := false
		for _, emp := range emps {
			if emp.EmployeeID == testEmp.EmployeeID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Employee not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testEmp.Surname,
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
		exists, err := repo.Exists(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Employee to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Employee to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondEmp := &domain.Employee{
			UserID:         nil,
			Surname:        "Сидоров",
			Name:           "Илья",
			FatherName:     nil,
			Birthday:       time.Date(1985, 12, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:    "79990001122",
			Job:            "Администратор",
			WorkExperience: 10,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondEmp)
		if err != nil {
			t.Fatalf("Create second Employee failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testEmp.EmployeeID, secondEmp.EmployeeID}
		emps, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(emps) != 2 {
			t.Errorf("Expected 2 Employees, got %d", len(emps))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Employee to be deleted")
		}
	})
}
=======
package db_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	_ "github.com/lib/pq"
)

func TestEmployeeRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewEmployeeRepository(sqlDB)
	FatherName := "Loxovich"

	testEmp := &domain.Employee{
		UserID:         nil,
		Surname:        "Иванов",
		Name:           "Иван",
		FatherName:     &FatherName,
		Birthday:       time.Date(1990, 5, 20, 0, 0, 0, 0, time.UTC),
		PhoneNumber:    "79991112233",
		Job:            "Программист",
		WorkExperience: 5,
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testEmp)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testEmp.EmployeeID == 0 {
			t.Error("Expected EmployeeID to be set after Create")
		} else {
			t.Logf("Created EmployeeID: %d", testEmp.EmployeeID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
		t.Logf("testEmp.EmployeeID = %d", testEmp.EmployeeID)
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Logf("Trying to get employee with ID: %d", testEmp.EmployeeID)
		emp, err := repo.GetByID(ctx, testEmp.EmployeeID)
		t.Logf("Successfully retrieved employee: %+v", emp)

		// Проверки полей
		if emp.Surname != testEmp.Surname {
			t.Errorf("Expected Surname %q, got %q", testEmp.Surname, emp.Surname)
		}
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if emp == nil {
			t.Fatal("Expected Employee to be found")
		}
		if emp.Surname != testEmp.Surname {
			t.Errorf("Expected Surname %q, got %q", testEmp.Surname, emp.Surname)
		}
		if emp.Name != testEmp.Name {
			t.Errorf("Expected Name %q, got %q", testEmp.Name, emp.Name)
		}
		if !emp.Birthday.Equal(testEmp.Birthday) {
			t.Errorf("Expected Birthday %v, got %v", testEmp.Birthday, emp.Birthday)
		}
		if emp.PhoneNumber != testEmp.PhoneNumber {
			t.Errorf("Expected PhoneNumber %q, got %q", testEmp.PhoneNumber, emp.PhoneNumber)
		}
		if emp.Job != testEmp.Job {
			t.Errorf("Expected Job %q, got %q", testEmp.Job, emp.Job)
		}
		if emp.WorkExperience != testEmp.WorkExperience {
			t.Errorf("Expected WorkExperience %d, got %d", testEmp.WorkExperience, emp.WorkExperience)
		}
		fmt.Printf("Original: %+v", testEmp)
		fmt.Printf("FromDB: %+v", emp)
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedEmp := *testEmp
		updatedEmp.Surname = "Петров"
		updatedEmp.Name = "Пётр"
		updatedEmp.PhoneNumber = "78882223344"
		updatedEmp.Job = "Тестировщик"
		updatedEmp.WorkExperience = 7

		err = repoWithTx.Update(ctx, &updatedEmp)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		emp, err := repo.GetByID(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if emp.Surname != "Петров" {
			t.Errorf("Expected Surname 'Петров' after update, got %q", emp.Surname)
		}
		if emp.Name != "Пётр" {
			t.Errorf("Expected Name 'Пётр' after update, got %q", emp.Name)
		}
		if emp.PhoneNumber != "78882223344" {
			t.Errorf("Expected PhoneNumber '78882223344' after update, got %q", emp.PhoneNumber)
		}
		if emp.Job != "Тестировщик" {
			t.Errorf("Expected Job 'Тестировщик' after update, got %q", emp.Job)
		}
		if emp.WorkExperience != 7 {
			t.Errorf("Expected WorkExperience 7 after update, got %d", emp.WorkExperience)
		}

		*testEmp = updatedEmp
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testEmp.Surname,
				},
			},
			Limit: 10,
		}

		emps, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(emps) == 0 {
			t.Error("Expected at least one Employee in List")
		} else {
			t.Logf("List returned %d items", len(emps))
		}

		found := false
		for _, emp := range emps {
			if emp.EmployeeID == testEmp.EmployeeID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Employee not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testEmp.Surname,
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
		exists, err := repo.Exists(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Employee to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Employee to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondEmp := &domain.Employee{
			UserID:         nil,
			Surname:        "Сидоров",
			Name:           "Илья",
			FatherName:     nil,
			Birthday:       time.Date(1985, 12, 15, 0, 0, 0, 0, time.UTC),
			PhoneNumber:    "79990001122",
			Job:            "Администратор",
			WorkExperience: 10,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondEmp)
		if err != nil {
			t.Fatalf("Create second Employee failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testEmp.EmployeeID, secondEmp.EmployeeID}
		emps, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(emps) != 2 {
			t.Errorf("Expected 2 Employees, got %d", len(emps))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testEmp.EmployeeID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Employee to be deleted")
		}
	})
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

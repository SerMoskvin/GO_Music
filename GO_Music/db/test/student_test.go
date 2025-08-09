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

func TestStudentRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentRepository(sqlDB)

	userID := 3
	fatherName := "Алексеевич"
	phoneNumber := "79109876654"

	testStudent := &domain.Student{
		UserID:        &userID,
		Surname:       "Иванов",
		Name:          "Иван",
		FatherName:    &fatherName,
		Birthday:      domain.ParseDMY("19.02.2000"),
		PhoneNumber:   &phoneNumber,
		GroupID:       1, // валидный GroupID
		MusprogrammID: 1, // валидный MusprogrammID
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testStudent)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testStudent.StudentID == 0 {
			t.Error("Expected StudentID to be set after Create")
		} else {
			t.Logf("Created StudentID: %d", testStudent.StudentID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		stud, err := repo.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if stud == nil {
			t.Fatal("Expected Student to be found")
		}
		if stud.Surname != testStudent.Surname {
			t.Errorf("Expected Surname %q, got %q", testStudent.Surname, stud.Surname)
		}
		if stud.Name != testStudent.Name {
			t.Errorf("Expected Name %q, got %q", testStudent.Name, stud.Name)
		}
		if stud.UserID == nil || *stud.UserID != *testStudent.UserID {
			t.Errorf("Expected UserID %v, got %v", testStudent.UserID, stud.UserID)
		}
		if stud.FatherName == nil || *stud.FatherName != *testStudent.FatherName {
			t.Errorf("Expected FatherName %v, got %v", testStudent.FatherName, stud.FatherName)
		}
		if !stud.Birthday.Equal(testStudent.Birthday) {
			t.Errorf("Expected Birthday %v, got %v", testStudent.Birthday, stud.Birthday)
		}
		if stud.PhoneNumber == nil || *stud.PhoneNumber != *testStudent.PhoneNumber {
			t.Errorf("Expected PhoneNumber %v, got %v", testStudent.PhoneNumber, stud.PhoneNumber)
		}
		if stud.GroupID != testStudent.GroupID {
			t.Errorf("Expected GroupID %d, got %d", testStudent.GroupID, stud.GroupID)
		}
		if stud.MusprogrammID != testStudent.MusprogrammID {
			t.Errorf("Expected MusprogrammID %d, got %d", testStudent.MusprogrammID, stud.MusprogrammID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedStudent := *testStudent
		updatedStudent.Surname = "Петров"
		updatedStudent.Name = "Пётр"
		newFatherName := "Сергеевич"
		updatedStudent.FatherName = &newFatherName
		newPhone := "10987654321"
		updatedStudent.PhoneNumber = &newPhone
		newUserID := 4
		updatedStudent.UserID = &newUserID
		updatedStudent.Birthday = updatedStudent.Birthday.AddDate(-1, 0, 0) // на год старше
		updatedStudent.GroupID = 2
		updatedStudent.MusprogrammID = 2

		err = repoWithTx.Update(ctx, &updatedStudent)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		stud, err := repo.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if stud.Surname != updatedStudent.Surname {
			t.Errorf("Expected Surname %q after update, got %q", updatedStudent.Surname, stud.Surname)
		}
		if stud.Name != updatedStudent.Name {
			t.Errorf("Expected Name %q after update, got %q", updatedStudent.Name, stud.Name)
		}
		if stud.FatherName == nil || *stud.FatherName != *updatedStudent.FatherName {
			t.Errorf("Expected FatherName %v after update, got %v", updatedStudent.FatherName, stud.FatherName)
		}
		if stud.PhoneNumber == nil || *stud.PhoneNumber != *updatedStudent.PhoneNumber {
			t.Errorf("Expected PhoneNumber %v after update, got %v", updatedStudent.PhoneNumber, stud.PhoneNumber)
		}
		if stud.UserID == nil || *stud.UserID != *updatedStudent.UserID {
			t.Errorf("Expected UserID %v after update, got %v", updatedStudent.UserID, stud.UserID)
		}
		if !stud.Birthday.Equal(updatedStudent.Birthday) {
			t.Errorf("Expected Birthday %v after update, got %v", updatedStudent.Birthday, stud.Birthday)
		}
		if stud.GroupID != updatedStudent.GroupID {
			t.Errorf("Expected GroupID %d after update, got %d", updatedStudent.GroupID, stud.GroupID)
		}
		if stud.MusprogrammID != updatedStudent.MusprogrammID {
			t.Errorf("Expected MusprogrammID %d after update, got %d", updatedStudent.MusprogrammID, stud.MusprogrammID)
		}

		*testStudent = updatedStudent
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testStudent.Surname,
				},
			},
			Limit: 10,
		}

		students, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(students) == 0 {
			t.Error("Expected at least one Student in List")
		} else {
			t.Logf("List returned %d items", len(students))
		}

		found := false
		for _, s := range students {
			if s.StudentID == testStudent.StudentID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Student not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testStudent.Surname,
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
		exists, err := repo.Exists(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Student to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Student to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondStudent := &domain.Student{
			UserID:        nil,
			Surname:       "Сидоров",
			Name:          "Сидор",
			FatherName:    nil,
			Birthday:      domain.ParseDMY("05.05.2005"),
			PhoneNumber:   nil,
			GroupID:       1,
			MusprogrammID: 1,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondStudent)
		if err != nil {
			t.Fatalf("Create second Student failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testStudent.StudentID, secondStudent.StudentID}
		students, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(students) != 2 {
			t.Errorf("Expected 2 Students, got %d", len(students))
		}
	})

	t.Run("SearchByName", func(t *testing.T) {
		results, err := repo.SearchByName(ctx, "Петр")
		if err != nil {
			t.Fatalf("SearchByName failed: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected at least one Student in SearchByName results")
		} else {
			t.Logf("SearchByName returned %d items", len(results))
		}

		found := false
		for _, s := range results {
			if s.StudentID == testStudent.StudentID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Student not found in SearchByName results")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Student to be deleted")
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

func TestStudentRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewStudentRepository(sqlDB)

	userID := 3
	fatherName := "Алексеевич"
	phoneNumber := "79109876654"

	testStudent := &domain.Student{
		UserID:        &userID,
		Surname:       "Иванов",
		Name:          "Иван",
		FatherName:    &fatherName,
		Birthday:      domain.ParseDMY("19.02.2000"),
		PhoneNumber:   &phoneNumber,
		GroupID:       1, // валидный GroupID
		MusprogrammID: 1, // валидный MusprogrammID
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testStudent)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testStudent.StudentID == 0 {
			t.Error("Expected StudentID to be set after Create")
		} else {
			t.Logf("Created StudentID: %d", testStudent.StudentID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		stud, err := repo.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if stud == nil {
			t.Fatal("Expected Student to be found")
		}
		if stud.Surname != testStudent.Surname {
			t.Errorf("Expected Surname %q, got %q", testStudent.Surname, stud.Surname)
		}
		if stud.Name != testStudent.Name {
			t.Errorf("Expected Name %q, got %q", testStudent.Name, stud.Name)
		}
		if stud.UserID == nil || *stud.UserID != *testStudent.UserID {
			t.Errorf("Expected UserID %v, got %v", testStudent.UserID, stud.UserID)
		}
		if stud.FatherName == nil || *stud.FatherName != *testStudent.FatherName {
			t.Errorf("Expected FatherName %v, got %v", testStudent.FatherName, stud.FatherName)
		}
		if !stud.Birthday.Equal(testStudent.Birthday) {
			t.Errorf("Expected Birthday %v, got %v", testStudent.Birthday, stud.Birthday)
		}
		if stud.PhoneNumber == nil || *stud.PhoneNumber != *testStudent.PhoneNumber {
			t.Errorf("Expected PhoneNumber %v, got %v", testStudent.PhoneNumber, stud.PhoneNumber)
		}
		if stud.GroupID != testStudent.GroupID {
			t.Errorf("Expected GroupID %d, got %d", testStudent.GroupID, stud.GroupID)
		}
		if stud.MusprogrammID != testStudent.MusprogrammID {
			t.Errorf("Expected MusprogrammID %d, got %d", testStudent.MusprogrammID, stud.MusprogrammID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedStudent := *testStudent
		updatedStudent.Surname = "Петров"
		updatedStudent.Name = "Пётр"
		newFatherName := "Сергеевич"
		updatedStudent.FatherName = &newFatherName
		newPhone := "10987654321"
		updatedStudent.PhoneNumber = &newPhone
		newUserID := 4
		updatedStudent.UserID = &newUserID
		updatedStudent.Birthday = updatedStudent.Birthday.AddDate(-1, 0, 0) // на год старше
		updatedStudent.GroupID = 2
		updatedStudent.MusprogrammID = 2

		err = repoWithTx.Update(ctx, &updatedStudent)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		stud, err := repo.GetByID(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if stud.Surname != updatedStudent.Surname {
			t.Errorf("Expected Surname %q after update, got %q", updatedStudent.Surname, stud.Surname)
		}
		if stud.Name != updatedStudent.Name {
			t.Errorf("Expected Name %q after update, got %q", updatedStudent.Name, stud.Name)
		}
		if stud.FatherName == nil || *stud.FatherName != *updatedStudent.FatherName {
			t.Errorf("Expected FatherName %v after update, got %v", updatedStudent.FatherName, stud.FatherName)
		}
		if stud.PhoneNumber == nil || *stud.PhoneNumber != *updatedStudent.PhoneNumber {
			t.Errorf("Expected PhoneNumber %v after update, got %v", updatedStudent.PhoneNumber, stud.PhoneNumber)
		}
		if stud.UserID == nil || *stud.UserID != *updatedStudent.UserID {
			t.Errorf("Expected UserID %v after update, got %v", updatedStudent.UserID, stud.UserID)
		}
		if !stud.Birthday.Equal(updatedStudent.Birthday) {
			t.Errorf("Expected Birthday %v after update, got %v", updatedStudent.Birthday, stud.Birthday)
		}
		if stud.GroupID != updatedStudent.GroupID {
			t.Errorf("Expected GroupID %d after update, got %d", updatedStudent.GroupID, stud.GroupID)
		}
		if stud.MusprogrammID != updatedStudent.MusprogrammID {
			t.Errorf("Expected MusprogrammID %d after update, got %d", updatedStudent.MusprogrammID, stud.MusprogrammID)
		}

		*testStudent = updatedStudent
	})

	t.Run("List", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testStudent.Surname,
				},
			},
			Limit: 10,
		}

		students, err := repo.List(ctx, filter)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(students) == 0 {
			t.Error("Expected at least one Student in List")
		} else {
			t.Logf("List returned %d items", len(students))
		}

		found := false
		for _, s := range students {
			if s.StudentID == testStudent.StudentID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Student not found in List")
		}
	})

	t.Run("Count", func(t *testing.T) {
		filter := db.Filter{
			Conditions: []db.Condition{
				{
					Field:    "surname",
					Operator: "=",
					Value:    testStudent.Surname,
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
		exists, err := repo.Exists(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected Student to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent Student to not exist")
		}
	})

	t.Run("GetByIDs", func(t *testing.T) {
		secondStudent := &domain.Student{
			UserID:        nil,
			Surname:       "Сидоров",
			Name:          "Сидор",
			FatherName:    nil,
			Birthday:      domain.ParseDMY("05.05.2005"),
			PhoneNumber:   nil,
			GroupID:       1,
			MusprogrammID: 1,
		}

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, secondStudent)
		if err != nil {
			t.Fatalf("Create second Student failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		ids := []int{testStudent.StudentID, secondStudent.StudentID}
		students, err := repo.GetByIDs(ctx, ids)
		if err != nil {
			t.Fatalf("GetByIDs failed: %v", err)
		}

		if len(students) != 2 {
			t.Errorf("Expected 2 Students, got %d", len(students))
		}
	})

	t.Run("SearchByName", func(t *testing.T) {
		results, err := repo.SearchByName(ctx, "Петр")
		if err != nil {
			t.Fatalf("SearchByName failed: %v", err)
		}

		if len(results) == 0 {
			t.Error("Expected at least one Student in SearchByName results")
		} else {
			t.Logf("SearchByName returned %d items", len(results))
		}

		found := false
		for _, s := range results {
			if s.StudentID == testStudent.StudentID {
				found = true
				break
			}
		}
		if !found {
			t.Error("Created Student not found in SearchByName results")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testStudent.StudentID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected Student to be deleted")
		}
	})
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

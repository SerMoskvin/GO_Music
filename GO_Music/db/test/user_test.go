<<<<<<< HEAD
package db_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	_ "github.com/lib/pq"
)

func TestUserRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewUserRepository(sqlDB)

	testUser := &domain.User{
		Login:            "testuser123",
		Password:         "password123",
		Role:             "user",
		Surname:          "Иванов",
		Name:             "Иван",
		RegistrationDate: time.Now().UTC().Truncate(time.Second),
		Email:            "ivanov@example.com",
		Image:            []byte{0xFF, 0xD8, 0xFF}, 
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testUser.UserID == 0 {
			t.Error("Expected UserID to be set after Create")
		} else {
			t.Logf("Created UserID: %d", testUser.UserID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		user, err := repo.GetByID(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if user == nil {
			t.Fatal("Expected User to be found")
		}
		if user.Login != testUser.Login {
			t.Errorf("Expected Login %q, got %q", testUser.Login, user.Login)
		}
		if user.Email != testUser.Email {
			t.Errorf("Expected Email %q, got %q", testUser.Email, user.Email)
		}
		if user.Surname != testUser.Surname {
			t.Errorf("Expected Surname %q, got %q", testUser.Surname, user.Surname)
		}
		if user.Name != testUser.Name {
			t.Errorf("Expected Name %q, got %q", testUser.Name, user.Name)
		}
		if !user.RegistrationDate.Equal(testUser.RegistrationDate) {
			t.Errorf("Expected RegistrationDate %v, got %v", testUser.RegistrationDate, user.RegistrationDate)
		}
		if len(user.Image) != len(testUser.Image) {
			t.Errorf("Expected Image length %d, got %d", len(testUser.Image), len(user.Image))
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedUser := *testUser
		updatedUser.Login = "updateduser"
		updatedUser.Email = "updated@example.com"
		updatedUser.Surname = "Петров"
		updatedUser.Name = "Пётр"
		updatedUser.Role = "admin"
		updatedUser.Image = []byte{0x89, 0x50, 0x4E, 0x47} // PNG header

		err = repoWithTx.Update(ctx, &updatedUser)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		user, err := repo.GetByID(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if user.Login != updatedUser.Login {
			t.Errorf("Expected Login %q after update, got %q", updatedUser.Login, user.Login)
		}
		if user.Email != updatedUser.Email {
			t.Errorf("Expected Email %q after update, got %q", updatedUser.Email, user.Email)
		}
		if user.Surname != updatedUser.Surname {
			t.Errorf("Expected Surname %q after update, got %q", updatedUser.Surname, user.Surname)
		}
		if user.Name != updatedUser.Name {
			t.Errorf("Expected Name %q after update, got %q", updatedUser.Name, user.Name)
		}
		if user.Role != updatedUser.Role {
			t.Errorf("Expected Role %q after update, got %q", updatedUser.Role, user.Role)
		}
		if len(user.Image) != len(updatedUser.Image) {
			t.Errorf("Expected Image length %d after update, got %d", len(updatedUser.Image), len(user.Image))
		}

		*testUser = updatedUser
	})

	t.Run("SearchByName", func(t *testing.T) {
		query := "Петро"
		users, err := repo.SearchByName(ctx, query)
		if err != nil {
			t.Fatalf("SearchByName failed: %v", err)
		}

		if len(users) == 0 {
			t.Error("Expected at least one user found by SearchByName")
		}

		found := false
		for _, u := range users {
			fullName := u.Surname + " " + u.Name
			if containsIgnoreCase(fullName, query) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("No user found with name containing %q", query)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected User to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent User to not exist")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected User to be deleted")
		}
	})
}

// containsIgnoreCase проверяет, содержится ли substr в s без учета регистра
func containsIgnoreCase(s, substr string) bool {
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}
=======
package db_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"GO_Music/config"
	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	_ "github.com/lib/pq"
)

func TestUserRepository_AllMethods(t *testing.T) {
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

	repo := repositories.NewUserRepository(sqlDB)

	testUser := &domain.User{
		Login:            "testuser123",
		Password:         "password123",
		Role:             "user",
		Surname:          "Иванов",
		Name:             "Иван",
		RegistrationDate: time.Now().UTC().Truncate(time.Second),
		Email:            "ivanov@example.com",
		Image:            []byte{0xFF, 0xD8, 0xFF}, // пример jpeg header
	}

	t.Run("Create", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if testUser.UserID == 0 {
			t.Error("Expected UserID to be set after Create")
		} else {
			t.Logf("Created UserID: %d", testUser.UserID)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		user, err := repo.GetByID(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("GetByID failed: %v", err)
		}
		if user == nil {
			t.Fatal("Expected User to be found")
		}
		if user.Login != testUser.Login {
			t.Errorf("Expected Login %q, got %q", testUser.Login, user.Login)
		}
		if user.Email != testUser.Email {
			t.Errorf("Expected Email %q, got %q", testUser.Email, user.Email)
		}
		if user.Surname != testUser.Surname {
			t.Errorf("Expected Surname %q, got %q", testUser.Surname, user.Surname)
		}
		if user.Name != testUser.Name {
			t.Errorf("Expected Name %q, got %q", testUser.Name, user.Name)
		}
		if !user.RegistrationDate.Equal(testUser.RegistrationDate) {
			t.Errorf("Expected RegistrationDate %v, got %v", testUser.RegistrationDate, user.RegistrationDate)
		}
		if len(user.Image) != len(testUser.Image) {
			t.Errorf("Expected Image length %d, got %d", len(testUser.Image), len(user.Image))
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		updatedUser := *testUser
		updatedUser.Login = "updateduser"
		updatedUser.Email = "updated@example.com"
		updatedUser.Surname = "Петров"
		updatedUser.Name = "Пётр"
		updatedUser.Role = "admin"
		updatedUser.Image = []byte{0x89, 0x50, 0x4E, 0x47} // PNG header

		err = repoWithTx.Update(ctx, &updatedUser)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		user, err := repo.GetByID(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("GetByID after update failed: %v", err)
		}
		if user.Login != updatedUser.Login {
			t.Errorf("Expected Login %q after update, got %q", updatedUser.Login, user.Login)
		}
		if user.Email != updatedUser.Email {
			t.Errorf("Expected Email %q after update, got %q", updatedUser.Email, user.Email)
		}
		if user.Surname != updatedUser.Surname {
			t.Errorf("Expected Surname %q after update, got %q", updatedUser.Surname, user.Surname)
		}
		if user.Name != updatedUser.Name {
			t.Errorf("Expected Name %q after update, got %q", updatedUser.Name, user.Name)
		}
		if user.Role != updatedUser.Role {
			t.Errorf("Expected Role %q after update, got %q", updatedUser.Role, user.Role)
		}
		if len(user.Image) != len(updatedUser.Image) {
			t.Errorf("Expected Image length %d after update, got %d", len(updatedUser.Image), len(user.Image))
		}

		*testUser = updatedUser
	})

	t.Run("SearchByName", func(t *testing.T) {
		// Ищем по фамилии+имени
		query := "Петро"
		users, err := repo.SearchByName(ctx, query)
		if err != nil {
			t.Fatalf("SearchByName failed: %v", err)
		}

		if len(users) == 0 {
			t.Error("Expected at least one user found by SearchByName")
		}

		found := false
		for _, u := range users {
			fullName := u.Surname + " " + u.Name
			if containsIgnoreCase(fullName, query) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("No user found with name containing %q", query)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected User to exist")
		}

		exists, err = repo.Exists(ctx, -1)
		if err != nil {
			t.Fatalf("Exists for non-existent id failed: %v", err)
		}
		if exists {
			t.Error("Expected non-existent User to not exist")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}
		defer tx.Rollback()

		repoWithTx := repo.WithTx(tx)

		err = repoWithTx.Delete(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("failed to commit tx: %v", err)
		}

		exists, err := repo.Exists(ctx, testUser.UserID)
		if err != nil {
			t.Fatalf("Exists after delete failed: %v", err)
		}
		if exists {
			t.Error("Expected User to be deleted")
		}
	})
}

// containsIgnoreCase проверяет, содержится ли substr в s без учета регистра
func containsIgnoreCase(s, substr string) bool {
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}
>>>>>>> 7267d2e1203c70e0401e4d5a7fe806cb4f2e2db7

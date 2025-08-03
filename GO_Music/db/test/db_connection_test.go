package db_test

import (
	"GO_Music/config"
	"GO_Music/db"
	"testing"
)

func TestInitPostgresDB(t *testing.T) {
	// [RU] Укажите путь к вашему реальному конфиг-файлу;
	// [ENG] Write path for your database config file
	cfgPath := "../../config/DB_config.yml"

	cfg, err := config.LoadDBConfig(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	sqlDB, err := db.InitPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	var result int
	err = sqlDB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}

	t.Log("Database connection test passed successfully")
}

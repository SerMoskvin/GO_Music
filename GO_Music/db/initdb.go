package db

import (
	"GO_Music/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// InitPostgresDB инициализирует подключение к базе данных PostgreSQL
func InitPostgresDB(cfg *config.DBConfig) (*sql.DB, error) {
	// Используем параметры из cfg.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

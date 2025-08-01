package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type ScheduleRepository struct {
	*postgreSQL.PostgresRepository[*domain.Schedule, int]
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.Schedule, int](
			db,
			"schedule",    // имя таблицы
			"schedule_id", // имя поля с ID
		),
	}
}

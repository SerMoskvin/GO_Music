package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type InstrumentRepository struct {
	*postgreSQL.PostgresRepository[domain.Instrument, int]
}

func NewInstrumentRepository(db *sql.DB) *InstrumentRepository {
	return &InstrumentRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.Instrument, int](
			db,
			"instrument",    // имя таблицы
			"instrument_id", // имя поля с ID
		),
	}
}

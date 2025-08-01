package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type ProgrammDistributionRepository struct {
	*postgreSQL.PostgresRepository[*domain.ProgrammDistribution, int]
}

func NewProgrammDistributionRepository(db *sql.DB) *ProgrammDistributionRepository {
	return &ProgrammDistributionRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.ProgrammDistribution, int](
			db,
			"programm_distribution", // имя таблицы
			"programm_distr_id",     // имя поля с ID
		),
	}
}

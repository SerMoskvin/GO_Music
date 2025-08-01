package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type SubjectDistributionRepository struct {
	*postgreSQL.PostgresRepository[*domain.SubjectDistribution, int]
}

func NewSubjectDistributionRepository(db *sql.DB) *SubjectDistributionRepository {
	return &SubjectDistributionRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.SubjectDistribution, int](
			db,
			"subject_distribution", // имя таблицы
			"subject_distr_id",     // имя поля с ID
		),
	}
}

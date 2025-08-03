package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type AudienceRepository struct {
	*postgreSQL.PostgresRepository[domain.Audience, int]
}

func NewAudienceRepository(db *sql.DB) *AudienceRepository {
	return &AudienceRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.Audience, int](
			db,
			"audience",    // имя таблицы
			"audience_id", // имя поля с ID
		),
	}
}

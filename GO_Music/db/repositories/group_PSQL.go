package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type StudyGroupRepository struct {
	*postgreSQL.PostgresRepository[domain.StudyGroup, int]
}

func NewStudyGroupRepository(db *sql.DB) *StudyGroupRepository {
	return &StudyGroupRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.StudyGroup, int](
			db,
			"study_group", // имя таблицы
			"group_id",    // имя поля с ID
		),
	}
}

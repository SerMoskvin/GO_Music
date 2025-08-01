package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type EmployeeRepository struct {
	*postgreSQL.PostgresRepository[*domain.Employee, int]
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.Employee, int](
			db,
			"employee",    // имя таблицы
			"employee_id", // имя поля с ID
		),
	}
}

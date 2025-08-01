package repositories

import (
	"context"
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type UserRepository struct {
	*postgreSQL.PostgresRepository[*domain.User, int]
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.User, int](
			db,
			"users",   // имя таблицы
			"user_id", // имя поля с ID
		),
	}
}

// Кастомные SQL-запросы для пользователей
const (
	searchUsersByNameQuery = `
		SELECT * FROM users 
		WHERE CONCAT(surname, ' ', name) ILIKE $1
		ORDER BY surname, name`
)

func (r *UserRepository) SearchByName(ctx context.Context, query string) ([]*domain.User, error) {
	rows, err := r.QueryContext(ctx, searchUsersByNameQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUserRows(rows)
}

func (r *UserRepository) scanUserRows(rows *sql.Rows) ([]*domain.User, error) {
	var users []*domain.User
	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.UserID,
			&u.Login,
			&u.Password,
			&u.Role,
			&u.Surname,
			&u.Name,
			&u.RegistrationDate,
			&u.Email,
			&u.Image,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

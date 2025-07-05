package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) engine.UserRepository {
	return &userRepository{db: db}
}

// Create вставляет нового пользователя в БД
func (r *userRepository) Create(user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	query := `
		INSERT INTO users 
			(login, password, role, surname, name, registration_date, image)
		VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7)
		RETURNING user_id
	`
	err := r.db.QueryRow(query,
		user.Login,
		user.Password,
		user.Role,
		user.Surname,
		user.Name,
		user.RegistrationDate,
		user.Image,
	).Scan(&user.UserID)
	if err != nil {
		return err
	}
	return nil
}

// Update обновляет данные пользователя в БД
func (r *userRepository) Update(user *domain.User) error {
	if user == nil || user.UserID == 0 {
		return errors.New("invalid user")
	}
	query := `
		UPDATE users SET
			login = \$1,
			password = \$2,
			role = \$3,
			surname = \$4,
			name = \$5,
			registration_date = \$6,
			image = \$7
		WHERE user_id = \$8
	`
	_, err := r.db.Exec(query,
		user.Login,
		user.Password,
		user.Role,
		user.Surname,
		user.Name,
		user.RegistrationDate,
		user.Image,
		user.UserID,
	)
	return err
}

// Delete удаляет пользователя по ID
func (r *userRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	query := `DELETE FROM users WHERE user_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetByID возвращает пользователя по ID
func (r *userRepository) GetByID(id int) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	query := `
		SELECT user_id, login, password, role, surname, name, registration_date, image
		FROM users WHERE user_id = \$1
	`
	row := r.db.QueryRow(query, id)
	user := &domain.User{}
	err := row.Scan(
		&user.UserID,
		&user.Login,
		&user.Password,
		&user.Role,
		&user.Surname,
		&user.Name,
		&user.RegistrationDate,
		&user.Image,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // пользователь не найден
		}
		return nil, err
	}
	return user, nil
}

// GetByIDs возвращает список пользователей по списку ID
func (r *userRepository) GetByIDs(ids []int) ([]*domain.User, error) {
	if len(ids) == 0 {
		return nil, errors.New("empty ids list")
	}

	// Формируем плейсхолдеры \$1, \$2, ...
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT user_id, login, password, role, surname, name, registration_date, image
		FROM users WHERE user_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.UserID,
			&user.Login,
			&user.Password,
			&user.Role,
			&user.Surname,
			&user.Name,
			&user.RegistrationDate,
			&user.Image,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

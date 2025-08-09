package repositories

import (
	"context"
	"database/sql"
	"time"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type StudentRepository struct {
	*postgreSQL.PostgresRepository[domain.Student, int]
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.Student, int](
			db,
			"student",    // имя таблицы
			"student_id", // имя поля с ID
		),
	}
}

const (
	searchStudentsByNameQuery = `
		SELECT * FROM student 
		WHERE CONCAT(surname, ' ', name, ' ', COALESCE(father_name, '')) ILIKE $1
		ORDER BY surname, name`
)

func (r *StudentRepository) SearchByName(ctx context.Context, query string) ([]*domain.Student, error) {
	rows, err := r.QueryContext(ctx, searchStudentsByNameQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanStudentRows(rows)
}

func (r *StudentRepository) scanStudentRows(rows *sql.Rows) ([]*domain.Student, error) {
	var students []*domain.Student
	for rows.Next() {
		var s domain.Student
		var fatherName, phoneNumber sql.NullString
		var userID sql.NullInt64
		var birthday time.Time

		err := rows.Scan(
			&s.StudentID,
			&userID,
			&s.Surname,
			&s.Name,
			&fatherName,
			&birthday,
			&phoneNumber,
			&s.GroupID,
			&s.MusprogrammID,
		)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			uid := int(userID.Int64)
			s.UserID = &uid
		}
		if fatherName.Valid {
			s.FatherName = &fatherName.String
		}
		if phoneNumber.Valid {
			s.PhoneNumber = &phoneNumber.String
		}
		s.Birthday = birthday

		students = append(students, &s)
	}
	return students, nil
}

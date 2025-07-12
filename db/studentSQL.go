package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type studentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) engine.StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(s *domain.Student) error {
	if s == nil {
		return errors.New("student is nil")
	}
	query := `
		INSERT INTO students
			(user_id, surname, name, father_name, birthday, phone_number, group_id, musprogramm_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING student_id
	`
	err := r.db.QueryRow(query,
		s.UserID,
		s.Surname,
		s.Name,
		s.FatherName,
		s.Birthday,
		s.PhoneNumber,
		s.GroupID,
		s.MusprogrammID,
	).Scan(&s.StudentID)
	return err
}

func (r *studentRepository) Update(s *domain.Student) error {
	if s == nil || s.StudentID == 0 {
		return errors.New("не указан ID студента")
	}
	query := `
        UPDATE students SET
            user_id = $1,
            surname = $2,
            name = $3,
            father_name = $4,
            birthday = $5,
            phone_number = $6,
            group_id = $7,
            musprogramm_id = $8
        WHERE student_id = $9
    `
	_, err := r.db.Exec(query,
		s.UserID,
		s.Surname,
		s.Name,
		s.FatherName,
		s.Birthday,
		s.PhoneNumber,
		s.GroupID,
		s.MusprogrammID,
		s.StudentID,
	)
	return err
}

func (r *studentRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID студента")
	}
	query := `DELETE FROM students WHERE student_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("студент не найден или уже удален")
	}
	return nil
}

func (r *studentRepository) GetByID(id int) (*domain.Student, error) {
	if id == 0 {
		return nil, errors.New("не указан ID студента")
	}
	query := `
	    SELECT student_id, user_id, surname, name, father_name, birthday, phone_number, group_id, musprogramm_id 
	    FROM students WHERE student_id = $1
    `
	row := r.db.QueryRow(query, id)

	var s domain.Student

	err := row.Scan(
		&s.StudentID,
		&s.UserID,
		&s.Surname,
		&s.Name,
		&s.FatherName,
		&s.Birthday,
		&s.PhoneNumber,
		&s.GroupID,
		&s.MusprogrammID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найдено
		}
		return nil, err
	}

	return &s, nil
}

func (r *studentRepository) GetByIDs(ids []int) ([]*domain.Student, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}

	var placeholders []string
	var args []interface{}
	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := fmt.Sprintf(`
        SELECT student_id, user_id, surname, name, father_name, birthday, phone_number, group_id, musprogramm_id 
        FROM students WHERE student_id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*domain.Student

	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(
			&s.StudentID,
			&s.UserID,
			&s.Surname,
			&s.Name,
			&s.FatherName,
			&s.Birthday,
			&s.PhoneNumber,
			&s.GroupID,
			&s.MusprogrammID,
		); err != nil {
			return nil, err
		}
		students = append(students, &s)
	}

	return students, nil
}

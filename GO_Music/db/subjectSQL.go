package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type subjectRepository struct {
	db *sql.DB
}

func NewSubjectRepository(db *sql.DB) engine.SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(s *domain.Subject) error {
	if s == nil {
		return errors.New("subject is nil")
	}
	query := `
		INSERT INTO subjects (subject_name, subject_type, short_desc)
		VALUES ($1, $2, $3)
		RETURNING subject_id
	`
	err := r.db.QueryRow(query, s.SubjectName, s.SubjectType, s.ShortDesc).Scan(&s.SubjectID)
	return err
}

func (r *subjectRepository) Update(s *domain.Subject) error {
	if s == nil || s.SubjectID == 0 {
		return errors.New("не указан ID предмета")
	}
	query := `
        UPDATE subjects SET
            subject_name = $1,
            subject_type = $2,
            short_desc = $3
        WHERE subject_id = $4
    `
	result, err := r.db.Exec(query, s.SubjectName, s.SubjectType, s.ShortDesc, s.SubjectID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("предмет не найден или уже обновлен")
	}
	return nil
}

func (r *subjectRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID предмета")
	}
	query := `DELETE FROM subjects WHERE subject_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("предмет не найден или уже удален")
	}
	return nil
}

func (r *subjectRepository) GetByID(id int) (*domain.Subject, error) {
	if id == 0 {
		return nil, errors.New("не указан ID предмета")
	}
	query := `SELECT subject_id, subject_name, subject_type, short_desc FROM subjects WHERE subject_id = $1`
	row := r.db.QueryRow(query, id)

	var s domain.Subject
	err := row.Scan(&s.SubjectID, &s.SubjectName, &s.SubjectType, &s.ShortDesc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найдено
		}
		return nil, err
	}
	return &s, nil
}

func (r *subjectRepository) GetByIDs(ids []int) ([]*domain.Subject, error) {
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
        SELECT subject_id, subject_name, subject_type, short_desc FROM subjects WHERE subject_id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*domain.Subject

	for rows.Next() {
		var s domain.Subject
		if err := rows.Scan(&s.SubjectID, &s.SubjectName, &s.SubjectType, &s.ShortDesc); err != nil {
			return nil, err
		}
		subjects = append(subjects, &s)
	}

	return subjects, nil
}

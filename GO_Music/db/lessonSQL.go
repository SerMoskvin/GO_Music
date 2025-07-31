package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
)

type lessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) engine.LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) Create(lesson *domain.Lesson) error {
	if lesson == nil {
		return errors.New("lesson is nil")
	}
	query := `
		INSERT INTO lessons
			(audience_id, employee_id, group_id, student_id, lesson_name, subject_id)
		VALUES (\$1, \$2, \$3, \$4, \$5, \$6)
		RETURNING lesson_id
	`
	err := r.db.QueryRow(query,
		lesson.AudienceID,
		lesson.EmployeeID,
		lesson.GroupID,
		lesson.StudentID,
		lesson.LessonName,
		lesson.SubjectID,
	).Scan(&lesson.LessonID)
	return err
}

func (r *lessonRepository) Update(lesson *domain.Lesson) error {
	if lesson == nil || lesson.LessonID == 0 {
		return errors.New("не указан ID урока")
	}
	query := `
		UPDATE lessons SET
			audience_id = \$1,
			employee_id = \$2,
			group_id = \$3,
			student_id = \$4,
			lesson_name = \$5,
			subject_id = \$6
		WHERE lesson_id = \$7
	`
	_, err := r.db.Exec(query,
		lesson.AudienceID,
		lesson.EmployeeID,
		lesson.GroupID,
		lesson.StudentID,
		lesson.LessonName,
		lesson.SubjectID,
		lesson.LessonID,
	)
	return err
}

func (r *lessonRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID урока")
	}
	query := `DELETE FROM lessons WHERE lesson_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *lessonRepository) GetByID(id int) (*domain.Lesson, error) {
	if id == 0 {
		return nil, errors.New("не указан ID урока")
	}
	query := `
		SELECT lesson_id, audience_id, employee_id, group_id, student_id, lesson_name, subject_id
		FROM lessons WHERE lesson_id = \$1
	`
	row := r.db.QueryRow(query, id)
	lesson := &domain.Lesson{}
	err := row.Scan(
		&lesson.LessonID,
		&lesson.AudienceID,
		&lesson.EmployeeID,
		&lesson.GroupID,
		&lesson.StudentID,
		&lesson.LessonName,
		&lesson.SubjectID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return lesson, nil
}

func (r *lessonRepository) GetByIDs(ids []int) ([]*domain.Lesson, error) {
	if len(ids) == 0 {
		return nil, errors.New("список ID пуст")
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT lesson_id, audience_id, employee_id, group_id, student_id, lesson_name, subject_id
		FROM lessons WHERE lesson_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		lesson := &domain.Lesson{}
		err := rows.Scan(
			&lesson.LessonID,
			&lesson.AudienceID,
			&lesson.EmployeeID,
			&lesson.GroupID,
			&lesson.StudentID,
			&lesson.LessonName,
			&lesson.SubjectID,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

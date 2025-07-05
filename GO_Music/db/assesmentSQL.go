package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type studentAssessmentRepository struct {
	db *sql.DB
}

func NewStudentAssessmentRepository(db *sql.DB) engine.StudentAssessmentRepository {
	return &studentAssessmentRepository{db: db}
}

func (r *studentAssessmentRepository) Create(sa *domain.StudentAssessment) error {
	if sa == nil {
		return errors.New("student assessment is nil")
	}
	query := `
		INSERT INTO student_assessment
			(lesson_id, student_id, task_type, grade, assessment_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING assessment_note_id
	`
	err := r.db.QueryRow(query,
		sa.LessonID,
		sa.StudentID,
		sa.TaskType,
		sa.Grade,
		sa.AssessmentDate,
	).Scan(&sa.AssessmentNoteID)
	return err
}

func (r *studentAssessmentRepository) Update(sa *domain.StudentAssessment) error {
	if sa == nil || sa.AssessmentNoteID == 0 {
		return errors.New("не указан ID оценки")
	}
	query := `
        UPDATE student_assessment SET
            lesson_id = $1,
            student_id = $2,
            task_type = $3,
            grade = $4,
            assessment_date = $5
        WHERE assessment_note_id = $6
    `
	_, err := r.db.Exec(query,
		sa.LessonID,
		sa.StudentID,
		sa.TaskType,
		sa.Grade,
		sa.AssessmentDate,
		sa.AssessmentNoteID,
	)
	return err
}

func (r *studentAssessmentRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID оценки")
	}
	query := `DELETE FROM student_assessment WHERE assessment_note_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("оценка не найдена или уже удалена")
	}
	return nil
}

func (r *studentAssessmentRepository) GetByID(id int) (*domain.StudentAssessment, error) {
	if id == 0 {
		return nil, errors.New("не указан ID оценки")
	}
	query := `
	    SELECT assessment_note_id, lesson_id, student_id, task_type, grade, assessment_date 
	    FROM student_assessment WHERE assessment_note_id = $1
    `
	row := r.db.QueryRow(query, id)

	var sa domain.StudentAssessment

	err := row.Scan(
		&sa.AssessmentNoteID,
		&sa.LessonID,
		&sa.StudentID,
		&sa.TaskType,
		&sa.Grade,
		&sa.AssessmentDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найдено
		}
		return nil, err
	}

	return &sa, nil
}

func (r *studentAssessmentRepository) GetByIDs(ids []int) ([]*domain.StudentAssessment, error) {
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
        SELECT assessment_note_id, lesson_id, student_id, task_type, grade, assessment_date 
        FROM student_assessment WHERE assessment_note_id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []*domain.StudentAssessment

	for rows.Next() {
		var sa domain.StudentAssessment
		if err := rows.Scan(
			&sa.AssessmentNoteID,
			&sa.LessonID,
			&sa.StudentID,
			&sa.TaskType,
			&sa.Grade,
			&sa.AssessmentDate,
		); err != nil {
			return nil, err
		}
		assessments = append(assessments, &sa)
	}

	return assessments, nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
)

type studentAttendanceRepository struct {
	db *sql.DB
}

func NewStudentAttendanceRepository(db *sql.DB) engine.StudentAttendanceRepository {
	return &studentAttendanceRepository{db: db}
}

func (r *studentAttendanceRepository) Create(sa *domain.StudentAttendance) error {
	if sa == nil {
		return errors.New("student attendance is nil")
	}
	query := `
		INSERT INTO student_attendance
			(student_id, lesson_id, presence_mark, attendance_date)
		VALUES ($1, $2, $3, $4)
		RETURNING attendance_note_id
	`
	err := r.db.QueryRow(query,
		sa.StudentID,
		sa.LessonID,
		sa.PresenceMark,
		sa.AttendanceDate,
	).Scan(&sa.AttendanceNoteID)
	return err
}

func (r *studentAttendanceRepository) Update(sa *domain.StudentAttendance) error {
	if sa == nil || sa.AttendanceNoteID == 0 {
		return errors.New("не указан ID записи посещения")
	}
	query := `
        UPDATE student_attendance SET
            student_id = $1,
            lesson_id = $2,
            presence_mark = $3,
            attendance_date = $4
        WHERE attendance_note_id = $5
    `
	_, err := r.db.Exec(query,
		sa.StudentID,
		sa.LessonID,
		sa.PresenceMark,
		sa.AttendanceDate,
		sa.AttendanceNoteID,
	)
	return err
}

func (r *studentAttendanceRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID записи посещения")
	}
	query := `DELETE FROM student_attendance WHERE attendance_note_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("запись посещения не найдена или уже удалена")
	}
	return nil
}

func (r *studentAttendanceRepository) GetByID(id int) (*domain.StudentAttendance, error) {
	if id == 0 {
		return nil, errors.New("не указан ID записи посещения")
	}
	query := `
	    SELECT attendance_note_id, student_id, lesson_id, presence_mark, attendance_date 
	    FROM student_attendance WHERE attendance_note_id = $1
    `
	row := r.db.QueryRow(query, id)

	var sa domain.StudentAttendance

	err := row.Scan(
		&sa.AttendanceNoteID,
		&sa.StudentID,
		&sa.LessonID,
		&sa.PresenceMark,
		&sa.AttendanceDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найдено
		}
		return nil, err
	}

	return &sa, nil
}

func (r *studentAttendanceRepository) GetByIDs(ids []int) ([]*domain.StudentAttendance, error) {
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
        SELECT attendance_note_id, student_id, lesson_id, presence_mark, attendance_date 
        FROM student_attendance WHERE attendance_note_id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*domain.StudentAttendance

	for rows.Next() {
		var sa domain.StudentAttendance
		if err := rows.Scan(
			&sa.AttendanceNoteID,
			&sa.StudentID,
			&sa.LessonID,
			&sa.PresenceMark,
			&sa.AttendanceDate,
		); err != nil {
			return nil, err
		}
		attendances = append(attendances, &sa)
	}

	return attendances, nil
}

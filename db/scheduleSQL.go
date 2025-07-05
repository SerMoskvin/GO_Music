package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type scheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) engine.ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(sch *domain.Schedule) error {
	if sch == nil {
		return errors.New("schedule is nil")
	}
	query := `
		INSERT INTO schedule
			(lesson_id, day_week, time_begin, time_end, sched_date_start, sched_date_end)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING schedule_id
	`
	err := r.db.QueryRow(query,
		sch.LessonID,
		sch.DayWeek,
		sch.TimeBegin,
		sch.TimeEnd,
		sch.SchedDateStart,
		sch.SchedDateEnd,
	).Scan(&sch.ScheduleID)
	return err
}

func (r *scheduleRepository) Update(sch *domain.Schedule) error {
	if sch == nil || sch.ScheduleID == 0 {
		return errors.New("не указан ID расписания")
	}
	query := `
        UPDATE schedule SET
            lesson_id = $1,
            day_week = $2,
            time_begin = $3,
            time_end = $4,
            sched_date_start = $5,
            sched_date_end = $6
        WHERE schedule_id = $7
    `
	_, err := r.db.Exec(query,
		sch.LessonID,
		sch.DayWeek,
		sch.TimeBegin,
		sch.TimeEnd,
		sch.SchedDateStart,
		sch.SchedDateEnd,
		sch.ScheduleID,
	)
	return err
}

func (r *scheduleRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID расписания")
	}
	query := `DELETE FROM schedule WHERE schedule_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("расписание не найдено или уже удалено")
	}
	return nil
}

func (r *scheduleRepository) GetByID(id int) (*domain.Schedule, error) {
	if id == 0 {
		return nil, errors.New("не указан ID расписания")
	}
	query := `
	    SELECT schedule_id, lesson_id, day_week, time_begin, time_end, sched_date_start, sched_date_end 
	    FROM schedule WHERE schedule_id = $1
    `
	row := r.db.QueryRow(query, id)

	var sch domain.Schedule

	err := row.Scan(
		&sch.ScheduleID,
		&sch.LessonID,
		&sch.DayWeek,
		&sch.TimeBegin,
		&sch.TimeEnd,
		&sch.SchedDateStart,
		&sch.SchedDateEnd,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // не найдено
		}
		return nil, err
	}

	return &sch, nil
}

func (r *scheduleRepository) GetByIDs(ids []int) ([]*domain.Schedule, error) {
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
        SELECT schedule_id, lesson_id, day_week, time_begin, time_end, sched_date_start, sched_date_end 
        FROM schedule WHERE schedule_id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*domain.Schedule

	for rows.Next() {
		var sch domain.Schedule
		if err := rows.Scan(
			&sch.ScheduleID,
			&sch.LessonID,
			&sch.DayWeek,
			&sch.TimeBegin,
			&sch.TimeEnd,
			&sch.SchedDateStart,
			&sch.SchedDateEnd,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, &sch)
	}

	return schedules, nil
}

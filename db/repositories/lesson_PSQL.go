package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type LessonRepository struct {
	*postgreSQL.PostgresRepository[domain.Lesson, int]
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.Lesson, int](
			db,
			"lesson",    // имя таблицы
			"lesson_id", // имя поля с ID
		),
		db: db,
	}
}

// CheckEmployeeAvailability проверяет, свободен ли преподаватель в указанное время
func (r *LessonRepository) CheckEmployeeAvailability(
	ctx context.Context,
	employeeID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	query := `
		SELECT NOT EXISTS(
			SELECT 1 FROM lesson l
			JOIN schedule s ON l.lesson_id = s.lesson_id
			WHERE l.employee_id = $1
			AND l.lesson_id != $4
			AND s.day_week = $5
			AND (
				(s.time_begin < $3 AND s.time_end > $2)
			)
		) AS is_available`

	var isAvailable bool
	dayWeek := startTime.Weekday().String()
	err := r.db.QueryRowContext(ctx, query,
		employeeID,
		startTime.Format("15:04"),
		endTime.Format("15:04"),
		excludeLessonID,
		dayWeek,
	).Scan(&isAvailable)

	if err != nil {
		return false, fmt.Errorf("failed to check employee availability: %w", err)
	}
	return isAvailable, nil
}

// CheckAudienceAvailability проверяет, свободна ли аудитория в указанное время
func (r *LessonRepository) CheckAudienceAvailability(
	ctx context.Context,
	audienceID int,
	startTime time.Time,
	endTime time.Time,
	excludeLessonID int,
) (bool, error) {
	query := `
		SELECT NOT EXISTS(
			SELECT 1 FROM lesson l
			JOIN schedule s ON l.lesson_id = s.lesson_id
			WHERE l.audience_id = $1
			AND l.lesson_id != $4
			AND s.day_week = $5
			AND (
				(s.time_begin < $3 AND s.time_end > $2)
			)
		) AS is_available`

	var isAvailable bool
	dayWeek := startTime.Weekday().String()
	err := r.db.QueryRowContext(ctx, query,
		audienceID,
		startTime.Format("15:04"),
		endTime.Format("15:04"),
		excludeLessonID,
		dayWeek,
	).Scan(&isAvailable)

	if err != nil {
		return false, fmt.Errorf("failed to check audience availability: %w", err)
	}
	return isAvailable, nil
}

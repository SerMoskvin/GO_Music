package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type StudentAttendanceRepository struct {
	*postgreSQL.PostgresRepository[*domain.StudentAttendance, int]
}

func NewStudentAttendanceRepository(db *sql.DB) *StudentAttendanceRepository {
	return &StudentAttendanceRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.StudentAttendance, int](
			db,
			"student_attendance", // имя таблицы
			"attendance_note_id", // имя поля с ID
		),
	}
}

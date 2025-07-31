package repositories

import (
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type StudentAssessmentRepository struct {
	*postgreSQL.PostgresRepository[*domain.StudentAssessment, int]
}

func NewStudentAssessmentRepository(db *sql.DB) *StudentAssessmentRepository {
	return &StudentAssessmentRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.StudentAssessment, int](
			db,
			"student_assessments", // имя таблицы
			"assessment_note_id",  // имя поля с ID
		),
	}
}

package repositories

import (
	"context"
	"database/sql"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type SubjectRepository struct {
	*postgreSQL.PostgresRepository[domain.Subject, int]
}

func NewSubjectRepository(db *sql.DB) *SubjectRepository {
	return &SubjectRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.Subject, int](
			db,
			"subject",    // имя таблицы
			"subject_id", // имя поля с ID
		),
	}
}

// Кастомные SQL-запросы для предметов
const (
	getPopularSubjectsQuery = `
		SELECT s.* FROM subject s
		JOIN programm_distribution pd ON s.subject_id = pd.subject_id
		GROUP BY s.subject_id
		ORDER BY COUNT(pd.musprogramm_id) DESC
		LIMIT $1`

	getSubjectsWithProgramsQuery = `
		SELECT s.* FROM subject s
		JOIN programm_distribution pd ON s.subject_id = pd.subject_id
		WHERE pd.musprogramm_id = $1
		ORDER BY s.subject_name`
)

func (r *SubjectRepository) GetPopularSubjects(ctx context.Context, limit int) ([]*domain.Subject, error) {
	rows, err := r.QueryContext(ctx, getPopularSubjectsQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjectRows(rows)
}

func (r *SubjectRepository) GetSubjectsWithPrograms(ctx context.Context, programID int) ([]*domain.Subject, error) {
	rows, err := r.QueryContext(ctx, getSubjectsWithProgramsQuery, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanSubjectRows(rows)
}

func (r *SubjectRepository) scanSubjectRows(rows *sql.Rows) ([]*domain.Subject, error) {
	var subjects []*domain.Subject
	for rows.Next() {
		var s domain.Subject
		err := rows.Scan(
			&s.SubjectID,
			&s.SubjectName,
			&s.SubjectType,
			&s.ShortDesc,
		)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, &s)
	}
	return subjects, nil
}

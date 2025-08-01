package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type ProgrammRepository struct {
	*postgreSQL.PostgresRepository[*domain.Programm, int]
	db *sql.DB
}

func NewProgrammRepository(db *sql.DB) *ProgrammRepository {
	return &ProgrammRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[*domain.Programm, int](
			db,
			"programm",       // имя таблицы
			"musprogramm_id", // имя поля с ID
		),
		db: db,
	}
}

// SearchByDescriptionFullText выполняет полнотекстовый поиск по описанию
func (r *ProgrammRepository) SearchByDescriptionFullText(ctx context.Context, searchText string) ([]*domain.Programm, error) {
	query := `
		SELECT * FROM programms 
		WHERE to_tsvector('russian', description) @@ to_tsquery('russian', $1)
		ORDER BY programm_name`

	rows, err := r.db.QueryContext(ctx, query, searchText)
	if err != nil {
		return nil, fmt.Errorf("failed to search by description: %w", err)
	}
	defer rows.Close()

	var programms []*domain.Programm
	for rows.Next() {
		var p domain.Programm
		if err := rows.Scan(
			&p.MusprogrammID,
			&p.ProgrammName,
			&p.ProgrammType,
			&p.Duration,
			&p.Instrument,
			&p.Description,
			&p.StudyLoad,
			&p.FinalCertificationForm,
		); err != nil {
			return nil, fmt.Errorf("failed to scan programm: %w", err)
		}
		programms = append(programms, &p)
	}

	return programms, nil
}

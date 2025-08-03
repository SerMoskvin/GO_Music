package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"GO_Music/db/postgreSQL"
	"GO_Music/domain"
)

type ProgrammRepository struct {
	*postgreSQL.PostgresRepository[domain.Programm, int]
	db *sql.DB
}

func NewProgrammRepository(db *sql.DB) *ProgrammRepository {
	return &ProgrammRepository{
		PostgresRepository: postgreSQL.NewPostgresRepository[domain.Programm, int](
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
		SELECT 
			musprogramm_id,
			programm_name,
			programm_type,
			duration,
			instrument,
			description,
			study_load,
			final_certification_form
		FROM programm 
		WHERE to_tsvector('russian', description) @@ plainto_tsquery('russian', $1)
		ORDER BY programm_name`

	rows, err := r.db.QueryContext(ctx, query, searchText)
	if err != nil {
		return nil, fmt.Errorf("failed to search by description: %w", err)
	}
	defer rows.Close()

	var programms []*domain.Programm
	for rows.Next() {
		var p domain.Programm
		var instrument sql.NullString
		var description sql.NullString

		err := rows.Scan(
			&p.MusprogrammID,
			&p.ProgrammName,
			&p.ProgrammType,
			&p.Duration,
			&instrument,
			&description,
			&p.StudyLoad,
			&p.FinalCertificationForm,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan programm: %w", err)
		}

		if instrument.Valid {
			p.Instrument = &instrument.String
		} else {
			p.Instrument = nil
		}
		if description.Valid {
			p.Description = &description.String
		} else {
			p.Description = nil
		}

		programms = append(programms, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return programms, nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
)

type programmRepository struct {
	db *sql.DB
}

func NewProgrammRepository(db *sql.DB) engine.ProgrammRepository {
	return &programmRepository{db: db}
}

func (r *programmRepository) Create(prog *domain.Programm) error {
	if prog == nil {
		return errors.New("programm is nil")
	}
	query := `
		INSERT INTO programm
			(programm_name, programm_type, duration, instrument, description, study_load, final_certification_form)
		VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7)
		RETURNING musprogramm_id
	`
	err := r.db.QueryRow(query,
		prog.ProgrammName,
		prog.ProgrammType,
		prog.Duration,
		prog.Instrument,
		prog.Description,
		prog.StudyLoad,
		prog.FinalCertificationForm,
	).Scan(&prog.MusprogrammID)
	return err
}

func (r *programmRepository) Update(prog *domain.Programm) error {
	if prog == nil || prog.MusprogrammID == 0 {
		return errors.New("не указан ID программы")
	}
	query := `
		UPDATE programm SET
			programm_name = \$1,
			programm_type = \$2,
			duration = \$3,
			instrument = \$4,
			description = \$5,
			study_load = \$6,
			final_certification_form = \$7
		WHERE musprogramm_id = \$8
	`
	_, err := r.db.Exec(query,
		prog.ProgrammName,
		prog.ProgrammType,
		prog.Duration,
		prog.Instrument,
		prog.Description,
		prog.StudyLoad,
		prog.FinalCertificationForm,
		prog.MusprogrammID,
	)
	return err
}

func (r *programmRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID программы")
	}
	query := `DELETE FROM programm WHERE musprogramm_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *programmRepository) GetByID(id int) (*domain.Programm, error) {
	if id == 0 {
		return nil, errors.New("не указан ID программы")
	}
	query := `
		SELECT musprogramm_id, programm_name, programm_type, duration, instrument, description, study_load, final_certification_form
		FROM programm
		WHERE musprogramm_id = \$1
	`
	prog := &domain.Programm{}
	err := r.db.QueryRow(query, id).Scan(
		&prog.MusprogrammID,
		&prog.ProgrammName,
		&prog.ProgrammType,
		&prog.Duration,
		&prog.Instrument,
		&prog.Description,
		&prog.StudyLoad,
		&prog.FinalCertificationForm,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return prog, nil
}

func (r *programmRepository) GetByIDs(ids []int) ([]*domain.Programm, error) {
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
		SELECT musprogramm_id, programm_name, programm_type, duration, instrument, description, study_load, final_certification_form
		FROM programm
		WHERE musprogramm_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []*domain.Programm
	for rows.Next() {
		prog := &domain.Programm{}
		err := rows.Scan(
			&prog.MusprogrammID,
			&prog.ProgrammName,
			&prog.ProgrammType,
			&prog.Duration,
			&prog.Instrument,
			&prog.Description,
			&prog.StudyLoad,
			&prog.FinalCertificationForm,
		)
		if err != nil {
			return nil, err
		}
		programs = append(programs, prog)
	}

	return programs, nil
}

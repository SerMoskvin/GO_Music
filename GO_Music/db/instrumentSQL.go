package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type instrumentRepository struct {
	db *sql.DB
}

func NewInstrumentRepository(db *sql.DB) engine.InstrumentRepository {
	return &instrumentRepository{db: db}
}

func (r *instrumentRepository) Create(instr *domain.Instrument) error {
	if instr == nil {
		return errors.New("инструмент не указан")
	}
	query := `
		INSERT INTO instruments
			(audience_id, name, instr_type, condition)
		VALUES (\$1, \$2, \$3, \$4)
		RETURNING instrument_id
	`
	err := r.db.QueryRow(query,
		instr.AudienceID,
		instr.Name,
		instr.InstrType,
		instr.Condition,
	).Scan(&instr.InstrumentID)
	return err
}

func (r *instrumentRepository) Update(instr *domain.Instrument) error {
	if instr == nil || instr.InstrumentID == 0 {
		return errors.New("не указан ID инструмента")
	}
	query := `
		UPDATE instruments SET
			audience_id = \$1,
			name = \$2,
			instr_type = \$3,
			condition = \$4
		WHERE instrument_id = \$5
	`
	_, err := r.db.Exec(query,
		instr.AudienceID,
		instr.Name,
		instr.InstrType,
		instr.Condition,
		instr.InstrumentID,
	)
	return err
}

func (r *instrumentRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID инструмента")
	}
	query := `DELETE FROM instruments WHERE instrument_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *instrumentRepository) GetByID(id int) (*domain.Instrument, error) {
	if id == 0 {
		return nil, errors.New("не указан ID инструмента")
	}
	query := `
		SELECT instrument_id, audience_id, name, instr_type, condition
		FROM instruments WHERE instrument_id = \$1
	`
	row := r.db.QueryRow(query, id)
	instr := &domain.Instrument{}
	err := row.Scan(
		&instr.InstrumentID,
		&instr.AudienceID,
		&instr.Name,
		&instr.InstrType,
		&instr.Condition,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return instr, nil
}

func (r *instrumentRepository) GetByIDs(ids []int) ([]*domain.Instrument, error) {
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
		SELECT instrument_id, audience_id, name, instr_type, condition
		FROM instruments WHERE instrument_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instruments []*domain.Instrument
	for rows.Next() {
		instr := &domain.Instrument{}
		err := rows.Scan(
			&instr.InstrumentID,
			&instr.AudienceID,
			&instr.Name,
			&instr.InstrType,
			&instr.Condition,
		)
		if err != nil {
			return nil, err
		}
		instruments = append(instruments, instr)
	}

	return instruments, nil
}

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
)

type audienceRepository struct {
	db *sql.DB
}

func NewAudienceRepository(db *sql.DB) engine.AudienceRepository {
	return &audienceRepository{db: db}
}

func (r *audienceRepository) Create(audience *domain.Audience) error {
	if audience == nil {
		return errors.New("audience is nil")
	}
	query := `
		INSERT INTO audiences 
			(name, audin_type, audin_number, capacity)
		VALUES (\$1, \$2, \$3, \$4)
		RETURNING audience_id
	`
	err := r.db.QueryRow(query,
		audience.Name,
		audience.AudinType,
		audience.AudinNumber,
		audience.Capacity,
	).Scan(&audience.AudienceID)
	if err != nil {
		return err
	}
	return nil
}

func (r *audienceRepository) Update(audience *domain.Audience) error {
	if audience == nil || audience.AudienceID == 0 {
		return errors.New("invalid audience")
	}
	query := `
		UPDATE audiences SET
			name = \$1,
			audin_type = \$2,
			audin_number = \$3,
			capacity = \$4
		WHERE audience_id = \$5
	`
	_, err := r.db.Exec(query,
		audience.Name,
		audience.AudinType,
		audience.AudinNumber,
		audience.Capacity,
		audience.AudienceID,
	)
	return err
}

func (r *audienceRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	query := `DELETE FROM audiences WHERE audience_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *audienceRepository) GetByID(id int) (*domain.Audience, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	query := `
		SELECT audience_id, name, audin_type, audin_number, capacity
		FROM audiences WHERE audience_id = \$1
	`
	row := r.db.QueryRow(query, id)
	audience := &domain.Audience{}
	err := row.Scan(
		&audience.AudienceID,
		&audience.Name,
		&audience.AudinType,
		&audience.AudinNumber,
		&audience.Capacity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return audience, nil
}

func (r *audienceRepository) GetByIDs(ids []int) ([]*domain.Audience, error) {
	if len(ids) == 0 {
		return nil, errors.New("empty ids list")
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT audience_id, name, audin_type, audin_number, capacity
		FROM audiences WHERE audience_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audiences []*domain.Audience
	for rows.Next() {
		audience := &domain.Audience{}
		err := rows.Scan(
			&audience.AudienceID,
			&audience.Name,
			&audience.AudinType,
			&audience.AudinNumber,
			&audience.Capacity,
		)
		if err != nil {
			return nil, err
		}
		audiences = append(audiences, audience)
	}

	return audiences, nil
}

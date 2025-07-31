package db

import (
	"database/sql"
	"errors"

	"GO_Music/domain"
)

type programmDistributionRepository struct {
	db *sql.DB
}

func NewProgrammDistributionRepository(db *sql.DB) engine.ProgrammDistributionRepository {
	return &programmDistributionRepository{db: db}
}

func (r *programmDistributionRepository) Create(pd *domain.ProgrammDistribution) error {
	if pd == nil {
		return errors.New("programm distribution is nil")
	}
	query := `
		INSERT INTO programm_distributions (musprogramm_id, subject_id)
		VALUES (\$1, \$2)
		RETURNING programm_distr_id
	`
	err := r.db.QueryRow(query, pd.MusprogrammID, pd.SubjectID).Scan(&pd.ProgrammDistrID)
	return err
}

func (r *programmDistributionRepository) Update(pd *domain.ProgrammDistribution) error {
	if pd == nil || pd.ProgrammDistrID == 0 {
		return errors.New("invalid programm_distr_id")
	}
	query := `
		UPDATE programm_distributions
		SET musprogramm_id = \$1, subject_id = \$2
		WHERE programm_distr_id = \$3
	`
	_, err := r.db.Exec(query, pd.MusprogrammID, pd.SubjectID, pd.ProgrammDistrID)
	return err
}

func (r *programmDistributionRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	query := `DELETE FROM programm_distributions WHERE programm_distr_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *programmDistributionRepository) GetByID(id int) (*domain.ProgrammDistribution, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	query := `
		SELECT programm_distr_id, musprogramm_id, subject_id
		FROM programm_distributions
		WHERE programm_distr_id = \$1
	`
	pd := &domain.ProgrammDistribution{}
	err := r.db.QueryRow(query, id).Scan(&pd.ProgrammDistrID, &pd.MusprogrammID, &pd.SubjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return pd, nil
}

func (r *programmDistributionRepository) ExistsByProgramAndSubject(musprogrammID, subjectID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM programm_distributions
			WHERE musprogramm_id = \$1 AND subject_id = \$2
		)
	`
	var exists bool
	err := r.db.QueryRow(query, musprogrammID, subjectID).Scan(&exists)
	return exists, err
}

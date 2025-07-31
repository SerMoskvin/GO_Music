package db

import (
	"database/sql"
	"errors"

	"GO_Music/domain"
)

type subjectDistributionRepository struct {
	db *sql.DB
}

func NewSubjectDistributionRepository(db *sql.DB) engine.SubjectDistributionRepository {
	return &subjectDistributionRepository{db: db}
}

func (r *subjectDistributionRepository) Create(sd *domain.SubjectDistribution) error {
	if sd == nil {
		return errors.New("subject distribution is nil")
	}
	query := `
		INSERT INTO subject_distributions (employee_id, subject_id)
		VALUES (\$1, \$2)
		RETURNING subject_distr_id
	`
	err := r.db.QueryRow(query, sd.EmployeeID, sd.SubjectID).Scan(&sd.SubjectDistrID)
	return err
}

func (r *subjectDistributionRepository) Update(sd *domain.SubjectDistribution) error {
	if sd == nil || sd.SubjectDistrID == 0 {
		return errors.New("invalid subject_distr_id")
	}
	query := `
		UPDATE subject_distributions
		SET employee_id = \$1, subject_id = \$2
		WHERE subject_distr_id = \$3
	`
	_, err := r.db.Exec(query, sd.EmployeeID, sd.SubjectID, sd.SubjectDistrID)
	return err
}

func (r *subjectDistributionRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	query := `DELETE FROM subject_distributions WHERE subject_distr_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *subjectDistributionRepository) GetByID(id int) (*domain.SubjectDistribution, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	query := `
		SELECT subject_distr_id, employee_id, subject_id
		FROM subject_distributions
		WHERE subject_distr_id = \$1
	`
	sd := &domain.SubjectDistribution{}
	err := r.db.QueryRow(query, id).Scan(&sd.SubjectDistrID, &sd.EmployeeID, &sd.SubjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return sd, nil
}

func (r *subjectDistributionRepository) ExistsByEmployeeAndSubject(employeeID, subjectID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM subject_distributions
			WHERE employee_id = \$1 AND subject_id = \$2
		)
	`
	var exists bool
	err := r.db.QueryRow(query, employeeID, subjectID).Scan(&exists)
	return exists, err
}

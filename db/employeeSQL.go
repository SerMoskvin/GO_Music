package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
	"GO_Music/engine"
)

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) engine.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	if employee == nil {
		return errors.New("employee is nil")
	}
	query := `
		INSERT INTO employees
			(user_id, surname, name, father_name, birthday, phone_number, job, work_experience)
		VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8)
		RETURNING employee_id
	`
	err := r.db.QueryRow(query,
		employee.UserID,
		employee.Surname,
		employee.Name,
		employee.FatherName,
		employee.Birthday,
		employee.PhoneNumber,
		employee.Job,
		employee.WorkExperience,
	).Scan(&employee.EmployeeID)
	if err != nil {
		return err
	}
	return nil
}

func (r *employeeRepository) Update(employee *domain.Employee) error {
	if employee == nil || employee.EmployeeID == 0 {
		return errors.New("invalid employee")
	}
	query := `
		UPDATE employees SET
			user_id = \$1,
			surname = \$2,
			name = \$3,
			father_name = \$4,
			birthday = \$5,
			phone_number = \$6,
			job = \$7,
			work_experience = \$8
		WHERE employee_id = \$9
	`
	_, err := r.db.Exec(query,
		employee.UserID,
		employee.Surname,
		employee.Name,
		employee.FatherName,
		employee.Birthday,
		employee.PhoneNumber,
		employee.Job,
		employee.WorkExperience,
		employee.EmployeeID,
	)
	return err
}

func (r *employeeRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	query := `DELETE FROM employees WHERE employee_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *employeeRepository) GetByID(id int) (*domain.Employee, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	query := `
		SELECT employee_id, user_id, surname, name, father_name, birthday, phone_number, job, work_experience
		FROM employees WHERE employee_id = \$1
	`
	row := r.db.QueryRow(query, id)
	employee := &domain.Employee{}
	err := row.Scan(
		&employee.EmployeeID,
		&employee.UserID,
		&employee.Surname,
		&employee.Name,
		&employee.FatherName,
		&employee.Birthday,
		&employee.PhoneNumber,
		&employee.Job,
		&employee.WorkExperience,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return employee, nil
}

func (r *employeeRepository) GetByIDs(ids []int) ([]*domain.Employee, error) {
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
		SELECT employee_id, user_id, surname, name, father_name, birthday, phone_number, job, work_experience
		FROM employees WHERE employee_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		employee := &domain.Employee{}
		err := rows.Scan(
			&employee.EmployeeID,
			&employee.UserID,
			&employee.Surname,
			&employee.Name,
			&employee.FatherName,
			&employee.Birthday,
			&employee.PhoneNumber,
			&employee.Job,
			&employee.WorkExperience,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

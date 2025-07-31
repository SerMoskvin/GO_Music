package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"GO_Music/domain"
)

type studyGroupRepository struct {
	db *sql.DB
}

func NewStudyGroupRepository(db *sql.DB) engine.StudyGroupRepository {
	return &studyGroupRepository{db: db}
}

func (r *studyGroupRepository) Create(group *domain.StudyGroup) error {
	if group == nil {
		return errors.New("группа обучения не указана")
	}
	query := `
		INSERT INTO study_groups
			(musprogramm_id, group_name, study_year, number_of_students)
		VALUES (\$1, \$2, \$3, \$4)
		RETURNING group_id
	`
	err := r.db.QueryRow(query,
		group.MusProgrammID,
		group.GroupName,
		group.StudyYear,
		group.NumberOfStudents,
	).Scan(&group.GroupID)
	return err
}

func (r *studyGroupRepository) Update(group *domain.StudyGroup) error {
	if group == nil || group.GroupID == 0 {
		return errors.New("не указан ID группы обучения")
	}
	query := `
		UPDATE study_groups SET
			musprogramm_id = \$1,
			group_name = \$2,
			study_year = \$3,
			number_of_students = \$4
		WHERE group_id = \$5
	`
	_, err := r.db.Exec(query,
		group.MusProgrammID,
		group.GroupName,
		group.StudyYear,
		group.NumberOfStudents,
		group.GroupID,
	)
	return err
}

func (r *studyGroupRepository) Delete(id int) error {
	if id == 0 {
		return errors.New("не указан ID группы обучения")
	}
	query := `DELETE FROM study_groups WHERE group_id = \$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *studyGroupRepository) GetByID(id int) (*domain.StudyGroup, error) {
	if id == 0 {
		return nil, errors.New("не указан ID группы обучения")
	}
	query := `
		SELECT group_id, musprogramm_id, group_name, study_year, number_of_students
		FROM study_groups WHERE group_id = \$1
	`
	row := r.db.QueryRow(query, id)
	group := &domain.StudyGroup{}
	err := row.Scan(
		&group.GroupID,
		&group.MusProgrammID,
		&group.GroupName,
		&group.StudyYear,
		&group.NumberOfStudents,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return group, nil
}

func (r *studyGroupRepository) GetByIDs(ids []int) ([]*domain.StudyGroup, error) {
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
		SELECT group_id, musprogramm_id, group_name, study_year, number_of_students
		FROM study_groups WHERE group_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.StudyGroup
	for rows.Next() {
		group := &domain.StudyGroup{}
		err := rows.Scan(
			&group.GroupID,
			&group.MusProgrammID,
			&group.GroupName,
			&group.StudyYear,
			&group.NumberOfStudents,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

package engine

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// ProgrammManager реализует бизнес-логику для музыкальных программ
type ProgrammManager struct {
	*BaseManager[int, *domain.Programm]
	db *sql.DB
}

func NewProgrammManager(
	repo db.Repository[*domain.Programm, int],
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ProgrammManager {
	return &ProgrammManager{
		BaseManager: NewBaseManager[int, *domain.Programm](repo, logger, txTimeout),
		db:          db,
	}
}

// GetByType возвращает программы указанного типа
func (m *ProgrammManager) GetByType(ctx context.Context, programmType string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "programm_type", Operator: "=", Value: programmType},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.logger.Error("GetByType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "type", Value: programmType},
		)
		return nil, fmt.Errorf("failed to get programms by type: %w", err)
	}
	return DereferenceSlice(programms), nil
}

// GetByInstrument возвращает программы для указанного инструмента
func (m *ProgrammManager) GetByInstrument(ctx context.Context, instrument string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "instrument", Operator: "=", Value: instrument},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.logger.Error("GetByInstrument failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return nil, fmt.Errorf("failed to get programms by instrument: %w", err)
	}
	return DereferenceSlice(programms), nil
}

// GetByName возвращает программу по точному названию
func (m *ProgrammManager) GetByName(ctx context.Context, name string) (*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "programm_name", Operator: "=", Value: name},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("GetByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return nil, fmt.Errorf("failed to get programm by name: %w", err)
	}

	if len(programms) == 0 {
		return nil, nil
	}
	return *programms[0], nil
}

// GetByDurationRange возвращает программы в указанном диапазоне длительности
func (m *ProgrammManager) GetByDurationRange(ctx context.Context, minDuration, maxDuration int) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "duration", Operator: ">=", Value: minDuration},
			{Field: "duration", Operator: "<=", Value: maxDuration},
		},
		OrderBy: "duration",
	})
	if err != nil {
		m.logger.Error("GetByDurationRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "min_duration", Value: minDuration},
			logger.Field{Key: "max_duration", Value: maxDuration},
		)
		return nil, fmt.Errorf("failed to get programms by duration range: %w", err)
	}
	return DereferenceSlice(programms), nil
}

// GetByStudyLoad возвращает программы с указанной учебной нагрузкой
func (m *ProgrammManager) GetByStudyLoad(ctx context.Context, studyLoad int) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "study_load", Operator: "=", Value: studyLoad},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.logger.Error("GetByStudyLoad failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "study_load", Value: studyLoad},
		)
		return nil, fmt.Errorf("failed to get programms by study load: %w", err)
	}
	return DereferenceSlice(programms), nil
}

// CheckNameUnique проверяет уникальность названия программы
func (m *ProgrammManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "programm_name", Operator: "=", Value: name},
			{Field: "musprogramm_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "exclude_id", Value: excludeID},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(programms) == 0, nil
}

// SearchByDescription возвращает программы, содержащие указанный текст в описании
func (m *ProgrammManager) SearchByDescription(ctx context.Context, searchText string) ([]*domain.Programm, error) {
	repo, ok := m.repo.(*repositories.ProgrammRepository)
	if !ok {
		// Если репозиторий не поддерживает полнотекстовый поиск, используем LIKE
		programms, err := m.List(ctx, db.Filter{
			Conditions: []db.Condition{
				{Field: "description", Operator: "LIKE", Value: "%" + searchText + "%"},
			},
			OrderBy: "programm_name",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search programms: %w", err)
		}
		return DereferenceSlice(programms), nil
	}

	return repo.SearchByDescriptionFullText(ctx, searchText)
}

// Create создает новую музыкальную программу
func (m *ProgrammManager) Create(ctx context.Context, programm *domain.Programm) error {
	if err := programm.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Проверяем уникальность названия
	isUnique, err := m.CheckNameUnique(ctx, programm.ProgrammName, 0)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("programm name %s already exists", programm.ProgrammName)
	}

	if err := m.repo.Create(ctx, &programm); err != nil {
		m.logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// Update обновляет данные музыкальной программы
func (m *ProgrammManager) Update(ctx context.Context, programm *domain.Programm) error {
	if err := programm.Validate(); err != nil {
		m.logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Проверяем уникальность названия
	isUnique, err := m.CheckNameUnique(ctx, programm.ProgrammName, programm.MusprogrammID)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("programm name %s already exists", programm.ProgrammName)
	}

	if err := m.repo.Update(ctx, &programm); err != nil {
		m.logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// BulkCreate массово создает программы в транзакции
func (m *ProgrammManager) BulkCreate(ctx context.Context, programms []*domain.Programm) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	for _, prog := range programms {
		if err := prog.Validate(); err != nil {
			return fmt.Errorf("validation failed for programm %v: %w", prog, err)
		}

		// Проверяем уникальность названия
		isUnique, err := m.CheckNameUnique(ctx, prog.ProgrammName, 0)
		if err != nil {
			return fmt.Errorf("name uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("programm name %s already exists", prog.ProgrammName)
		}

		ptrToProg := &prog
		if err := txRepo.Create(ctx, ptrToProg); err != nil {
			return fmt.Errorf("create failed for programm %v: %w", prog, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

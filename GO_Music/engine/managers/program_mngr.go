package managers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"GO_Music/db"
	"GO_Music/db/repositories"
	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
)

// ProgrammManager реализует бизнес-логику для музыкальных программ
type ProgrammManager struct {
	*engine.BaseManager[int, domain.Programm, *domain.Programm]
	repo *repositories.ProgrammRepository
	db   *sql.DB
}

func NewProgrammManager(
	repo *repositories.ProgrammRepository,
	db *sql.DB,
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ProgrammManager {
	return &ProgrammManager{
		BaseManager: engine.NewBaseManager[int, domain.Programm, *domain.Programm](repo, logger, txTimeout),
		repo:        repo,
		db:          db,
	}
}

// [RU] GetByType возвращает программы указанного типа <--->
// [ENG] GetByType returns programs of the specified type
func (m *ProgrammManager) GetByType(ctx context.Context, programmType string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "programm_type", Operator: "=", Value: programmType},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.Logger.Error("GetByType failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "type", Value: programmType},
		)
		return nil, fmt.Errorf("failed to get programms by type: %w", err)
	}
	return programms, nil
}

// [RU] GetByInstrument возвращает программы для указанного инструмента <--->
// [ENG] GetByInstrument returns programs for the specified instrument
func (m *ProgrammManager) GetByInstrument(ctx context.Context, instrument string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "instrument", Operator: "=", Value: instrument},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.Logger.Error("GetByInstrument failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "instrument", Value: instrument},
		)
		return nil, fmt.Errorf("failed to get programms by instrument: %w", err)
	}
	return programms, nil
}

// [RU] GetByName возвращает программу по точному названию <--->
// [ENG] GetByName returns a program by exact name
func (m *ProgrammManager) GetByName(ctx context.Context, name string) (*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "programm_name", Operator: "=", Value: name},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("GetByName failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return nil, fmt.Errorf("failed to get programm by name: %w", err)
	}
	if len(programms) == 0 {
		return nil, nil
	}
	return programms[0], nil
}

// [RU] GetByDurationRange возвращает программы в указанном диапазоне длительности <--->
// [ENG] GetByDurationRange returns programs in the specified duration range
func (m *ProgrammManager) GetByDurationRange(ctx context.Context, minDuration, maxDuration int) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "duration", Operator: ">=", Value: minDuration},
			{Field: "duration", Operator: "<=", Value: maxDuration},
		},
		OrderBy: "duration",
	})
	if err != nil {
		m.Logger.Error("GetByDurationRange failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "min_duration", Value: minDuration},
			logger.Field{Key: "max_duration", Value: maxDuration},
		)
		return nil, fmt.Errorf("failed to get programms by duration range: %w", err)
	}
	return programms, nil
}

// [RU] GetByStudyLoad возвращает программы с указанной учебной нагрузкой <--->
// [ENG] GetByStudyLoad returns programs with the specified study load
func (m *ProgrammManager) GetByStudyLoad(ctx context.Context, studyLoad int) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "study_load", Operator: "=", Value: studyLoad},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.Logger.Error("GetByStudyLoad failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "study_load", Value: studyLoad},
		)
		return nil, fmt.Errorf("failed to get programms by study load: %w", err)
	}
	return programms, nil
}

// [RU] CheckNameUnique проверяет уникальность названия программы <--->
// [ENG] CheckNameUnique checks the uniqueness of the program name
func (m *ProgrammManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	programms, err := m.List(ctx, db.Filter{
		Conditions: []db.Condition{
			{Field: "programm_name", Operator: "=", Value: name},
			{Field: "musprogramm_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.Logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "exclude_id", Value: excludeID},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(programms) == 0, nil
}

// [RU] SearchByDescription возвращает программы, содержащие указанный текст в описании <--->
// [ENG] SearchByDescription returns programs containing the specified text in the description
func (m *ProgrammManager) SearchByDescription(ctx context.Context, searchText string) ([]*domain.Programm, error) {
	programms, err := m.repo.SearchByDescriptionFullText(ctx, searchText)
	if err != nil {
		m.Logger.Error("SearchByDescription failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "search_text", Value: searchText},
		)
		return nil, fmt.Errorf("failed to search programms: %w", err)
	}
	return programms, nil
}

// [RU] Create создает новую музыкальную программу <--->
// [ENG] Create creates a new music program
func (m *ProgrammManager) Create(ctx context.Context, programm *domain.Programm) error {
	if err := programm.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, programm.ProgrammName, 0)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("programm name %s already exists", programm.ProgrammName)
	}

	if err := m.Repo.Create(ctx, programm); err != nil {
		m.Logger.Error("Create failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}

// [RU] Update обновляет данные музыкальной программы <--->
// [ENG] Update updates the data of the music program
func (m *ProgrammManager) Update(ctx context.Context, programm *domain.Programm) error {
	if err := programm.Validate(); err != nil {
		m.Logger.Error("Validation failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("validation failed: %w", err)
	}

	isUnique, err := m.CheckNameUnique(ctx, programm.ProgrammName, programm.MusprogrammID)
	if err != nil {
		return fmt.Errorf("name uniqueness check failed: %w", err)
	}
	if !isUnique {
		return fmt.Errorf("programm name %s already exists", programm.ProgrammName)
	}

	if err := m.Repo.Update(ctx, programm); err != nil {
		m.Logger.Error("Update failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "programm", Value: programm},
		)
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

// [RU] BulkCreate массово создает программы в транзакции <--->
// [ENG] BulkCreate creates multiple programs in a transaction
func (m *ProgrammManager) BulkCreate(ctx context.Context, programms []*domain.Programm) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txRepo := m.Repo.WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, prog := range programms {
		if err := prog.Validate(); err != nil {
			return fmt.Errorf("validation failed for programm %v: %w", prog, err)
		}

		isUnique, err := m.CheckNameUnique(ctx, prog.ProgrammName, 0)
		if err != nil {
			return fmt.Errorf("name uniqueness check failed: %w", err)
		}
		if !isUnique {
			return fmt.Errorf("programm name %s already exists", prog.ProgrammName)
		}

		if err := txRepo.Create(ctx, prog); err != nil {
			return fmt.Errorf("create failed for programm %v: %w", prog, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

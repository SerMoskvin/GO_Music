package engine

import (
	"context"
	"fmt"
	"time"

	"GO_Music/domain"

	"github.com/SerMoskvin/logger"
)

// ProgrammManager реализует бизнес-логику для музыкальных программ
type ProgrammManager struct {
	*BaseManager[domain.Programm, *domain.Programm]
}

func NewProgrammManager(
	repo Repository[domain.Programm, *domain.Programm],
	logger *logger.LevelLogger,
	txTimeout time.Duration,
) *ProgrammManager {
	return &ProgrammManager{
		BaseManager: NewBaseManager[domain.Programm](repo, logger, txTimeout),
	}
}

// GetByType возвращает программы указанного типа
func (m *ProgrammManager) GetByType(ctx context.Context, programmType string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return programms, nil
}

// GetByInstrument возвращает программы для указанного инструмента
func (m *ProgrammManager) GetByInstrument(ctx context.Context, instrument string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return programms, nil
}

// GetByName возвращает программу по точному названию
func (m *ProgrammManager) GetByName(ctx context.Context, name string) (*domain.Programm, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return programms[0], nil
}

// GetByDurationRange возвращает программы в указанном диапазоне длительности
func (m *ProgrammManager) GetByDurationRange(ctx context.Context, minDuration, maxDuration int) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return programms, nil
}

// GetByStudyLoad возвращает программы с указанной учебной нагрузкой
func (m *ProgrammManager) GetByStudyLoad(ctx context.Context, studyLoad int) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
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
	return programms, nil
}

// CheckNameUnique проверяет уникальность названия программы
func (m *ProgrammManager) CheckNameUnique(ctx context.Context, name string, excludeID int) (bool, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "programm_name", Operator: "=", Value: name},
			{Field: "musprogramm_id", Operator: "!=", Value: excludeID},
		},
		Limit: 1,
	})
	if err != nil {
		m.logger.Error("CheckNameUnique failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name},
		)
		return false, fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	return len(programms) == 0, nil
}

// SearchByDescription возвращает программы, содержащие указанный текст в описании
func (m *ProgrammManager) SearchByDescription(ctx context.Context, searchText string) ([]*domain.Programm, error) {
	programms, err := m.List(ctx, Filter{
		Conditions: []Condition{
			{Field: "description", Operator: "LIKE", Value: "%" + searchText + "%"},
		},
		OrderBy: "programm_name",
	})
	if err != nil {
		m.logger.Error("SearchByDescription failed",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "search_text", Value: searchText},
		)
		return nil, fmt.Errorf("failed to search programms by description: %w", err)
	}
	return programms, nil
}

package engine

import (
	"time"

	"github.com/SerMoskvin/logger"
)

// Managers содержит ВСЕ менеджеры приложения
type Managers struct {
	Audience      *AudienceManager
	Employee      *EmployeeManager
	StudyGroup    *StudyGroupManager
	Instrument    *InstrumentManager
	ProgrammDistr *ProgrammDistributionManager
	SubjectDistr  *SubjectDistributionManager
}

// NewManagers создает ВСЕ менеджеры
func NewManagers(repos *Repositories, logger *logger.LevelLogger) *Managers {
	txTimeout := 10 * time.Second // Общий таймаут для всех менеджеров

	return &Managers{
		Audience:      NewAudienceManager(repos.Audience, logger, txTimeout),
		Employee:      NewEmployeeManager(repos.Employee, logger, txTimeout),
		StudyGroup:    NewStudyGroupManager(repos.StudyGroup, logger, txTimeout),
		Instrument:    NewInstrumentManager(repos.Instrument, logger, txTimeout),
		ProgrammDistr: NewProgrammDistributionManager(repos.ProgrammDistr, logger),
		SubjectDistr:  NewSubjectDistributionManager(repos.SubjectDistr, logger),
	}
}

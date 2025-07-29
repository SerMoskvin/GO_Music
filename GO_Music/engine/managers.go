package engine

import (
	"time"

	"github.com/SerMoskvin/access"
	"github.com/SerMoskvin/logger"
)

// Managers содержит все менеджеры приложения
type Managers struct {
	Assesment     *StudentAssessmentManager
	Attendance    *StudentAttendanceManager
	Audience      *AudienceManager
	Employee      *EmployeeManager
	StudyGroup    *StudyGroupManager
	Schedule      *ScheduleManager
	Instrument    *InstrumentManager
	ProgrammDistr *ProgrammDistributionManager
	SubjectDistr  *SubjectDistributionManager
	Lesson        *LessonManager
	Programm      *ProgrammManager
	Student       *StudentManager
	Subject       *SubjectManager
	User          *UserManager
}

// NewManagers создает все менеджеры
func NewManagers(repos *Repositories, logger *logger.LevelLogger, auth *access.Authenticator) *Managers {
	txTimeout := 10 * time.Second // Общий таймаут для всех менеджеров

	return &Managers{
		Audience:      NewAudienceManager(repos.Audience, logger, txTimeout),
		Assesment:     NewStudentAssessmentManager(repos.Assesment, logger, txTimeout),
		Attendance:    NewStudentAttendanceManager(repos.Attendance, logger, txTimeout),
		Employee:      NewEmployeeManager(repos.Employee, logger, txTimeout),
		StudyGroup:    NewStudyGroupManager(repos.StudyGroup, logger, txTimeout),
		Schedule:      NewScheduleManager(repos.Schedule, logger, txTimeout),
		Instrument:    NewInstrumentManager(repos.Instrument, logger, txTimeout),
		ProgrammDistr: NewProgrammDistributionManager(repos.ProgrammDistr, logger),
		SubjectDistr:  NewSubjectDistributionManager(repos.SubjectDistr, logger),
		Lesson:        NewLessonManager(repos.Lesson, logger, txTimeout),
		Programm:      NewProgrammManager(repos.Programm, logger, txTimeout),
		Student:       NewStudentManager(repos.Student, logger, txTimeout),
		Subject:       NewSubjectManager(repos.Subject, logger, txTimeout),
		User:          NewUserManager(repos.User, logger, txTimeout, auth),
	}
}

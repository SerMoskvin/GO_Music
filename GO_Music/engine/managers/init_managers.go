package managers

import (
	"database/sql"
	"time"

	"GO_Music/db/repositories"

	"github.com/SerMoskvin/access"
	"github.com/SerMoskvin/logger"
)

// Managers содержит все менеджеры приложения
type Managers struct {
	Assessment    *StudentAssessmentManager
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
func NewManagers(db *sql.DB, repos *repositories.Repositories, logger *logger.LevelLogger, auth *access.Authenticator) *Managers {
	txTimeout := 10 * time.Second // Общий таймаут для всех менеджеров

	return &Managers{
		Assessment:    NewStudentAssessmentManager(repos.Assessment, db, logger, txTimeout),
		Attendance:    NewStudentAttendanceManager(repos.Attendance, db, logger, txTimeout),
		Audience:      NewAudienceManager(repos.Audience, logger, txTimeout),
		Employee:      NewEmployeeManager(repos.Employee, db, logger, txTimeout),
		StudyGroup:    NewStudyGroupManager(repos.StudyGroup, db, logger, txTimeout),
		Schedule:      NewScheduleManager(repos.Schedule, db, logger, txTimeout),
		Instrument:    NewInstrumentManager(repos.Instrument, db, logger, txTimeout),
		ProgrammDistr: NewProgrammDistributionManager(repos.ProgrammDistr, db, logger, txTimeout),
		SubjectDistr:  NewSubjectDistributionManager(repos.SubjectDistr, db, logger, txTimeout),
		Lesson:        NewLessonManager(repos.Lesson, db, logger, txTimeout),
		Programm:      NewProgrammManager(repos.Programm, db, logger, txTimeout),
		Student:       NewStudentManager(repos.Student, db, logger, txTimeout),
		Subject:       NewSubjectManager(repos.Subject, db, logger, txTimeout),
		User:          NewUserManager(repos.User, db, logger, txTimeout, auth),
	}
}

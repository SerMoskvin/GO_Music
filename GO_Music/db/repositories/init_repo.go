package repositories

import (
	"database/sql"
)

// Repositories содержит все репозитории приложения
type Repositories struct {
	Audience      *AudienceRepository
	Assessment    *StudentAssessmentRepository
	Attendance    *StudentAttendanceRepository
	Employee      *EmployeeRepository
	StudyGroup    *StudyGroupRepository
	Schedule      *ScheduleRepository
	Instrument    *InstrumentRepository
	ProgrammDistr *ProgrammDistributionRepository
	SubjectDistr  *SubjectDistributionRepository
	Lesson        *LessonRepository
	Programm      *ProgrammRepository
	Student       *StudentRepository
	Subject       *SubjectRepository
	User          *UserRepository
}

// NewRepositories создает все репозитории
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Audience:      NewAudienceRepository(db),
		Assessment:    NewStudentAssessmentRepository(db),
		Attendance:    NewStudentAttendanceRepository(db),
		Employee:      NewEmployeeRepository(db),
		StudyGroup:    NewStudyGroupRepository(db),
		Schedule:      NewScheduleRepository(db),
		Instrument:    NewInstrumentRepository(db),
		ProgrammDistr: NewProgrammDistributionRepository(db),
		SubjectDistr:  NewSubjectDistributionRepository(db),
		Lesson:        NewLessonRepository(db),
		Programm:      NewProgrammRepository(db),
		Student:       NewStudentRepository(db),
		Subject:       NewSubjectRepository(db),
		User:          NewUserRepository(db),
	}
}

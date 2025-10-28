package handlers

import (
	"GO_Music/engine/managers"

	"github.com/SerMoskvin/logger"
	"github.com/go-chi/chi/v5"
)

// Handlers содержит все хендлеры приложения
type Handlers struct {
	Assessment    *StudentAssessmentHandler
	Attendance    *StudentAttendanceHandler
	Audience      *AudienceHandler
	Employee      *EmployeeHandler
	StudyGroup    *StudyGroupHandler
	Schedule      *ScheduleHandler
	Instrument    *InstrumentHandler
	ProgrammDistr *ProgrammDistributionHandler
	SubjectDistr  *SubjectDistributionHandler
	Lesson        *LessonHandler
	Programm      *ProgrammHandler
	Student       *StudentHandler
	Subject       *SubjectHandler
	User          *UserHandler
}

// NewHandlers создает все хендлеры
func NewHandlers(managers *managers.Managers, logger *logger.LevelLogger) *Handlers {
	return &Handlers{
		Assessment:    NewStudentAssessmentHandler(managers.Assessment, logger),
		Attendance:    NewStudentAttendanceHandler(managers.Attendance, logger),
		Audience:      NewAudienceHandler(managers.Audience, logger),
		Employee:      NewEmployeeHandler(managers.Employee, logger),
		StudyGroup:    NewStudyGroupHandler(managers.StudyGroup, logger),
		Schedule:      NewScheduleHandler(managers.Schedule, logger),
		Instrument:    NewInstrumentHandler(managers.Instrument, logger),
		ProgrammDistr: NewProgrammDistributionHandler(managers.ProgrammDistr, logger),
		SubjectDistr:  NewSubjectDistributionHandler(managers.SubjectDistr, logger),
		Lesson:        NewLessonHandler(managers.Lesson, logger),
		Programm:      NewProgrammHandler(managers.Programm, logger),
		Student:       NewStudentHandler(managers.Student, logger),
		Subject:       NewSubjectHandler(managers.Subject, logger),
		User:          NewUserHandler(managers.User, logger),
	}
}

func (h *Handlers) ToMap() map[string]interface{ Routes() chi.Router } {
	return map[string]interface{ Routes() chi.Router }{
		"assessments":            h.Assessment,
		"attendances":            h.Attendance,
		"audiences":              h.Audience,
		"employees":              h.Employee,
		"study-groups":           h.StudyGroup,
		"schedules":              h.Schedule,
		"instruments":            h.Instrument,
		"programm-distributions": h.ProgrammDistr,
		"subject-distributions":  h.SubjectDistr,
		"lessons":                h.Lesson,
		"programms":              h.Programm,
		"students":               h.Student,
		"subjects":               h.Subject,
		"users":                  h.User,
	}
}

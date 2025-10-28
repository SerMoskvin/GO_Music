package handlers

import (
	"GO_Music/api"
	dto "GO_Music/api/DTO"
	"GO_Music/domain"
	m "GO_Music/engine/managers"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SerMoskvin/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type LessonHandler struct {
	*api.BaseHandler[int, domain.Lesson, *domain.Lesson,
		dto.LessonCreateDTO, dto.LessonUpdateDTO, dto.LessonResponseDTO]
	manager *m.LessonManager
	mapper  *dto.LessonMapper
}

func NewLessonHandler(
	manager *m.LessonManager,
	logger *logger.LevelLogger,
) *LessonHandler {
	mapper := dto.NewLessonMapper()

	return &LessonHandler{
		BaseHandler: api.NewBaseHandler(
			manager.BaseManager,
			logger,
			mapper.ToDomain,
			mapper.UpdateDomain,
			mapper.ToResponse,
			nil,
			api.BaseHandlerConfig{
				DefaultPageSize: 20,
				MaxPageSize:     100,
			},
		),
		manager: manager,
		mapper:  mapper,
	}
}

func (h *LessonHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	r.Get("/by-employee/{employee_id}", h.GetByEmployee)
	r.Get("/by-group/{group_id}", h.GetByGroup)
	r.Get("/by-student/{student_id}", h.GetByStudent)
	r.Get("/by-subject/{subject_id}", h.GetBySubject)
	r.Get("/by-audience/{audience_id}", h.GetByAudience)
	r.Get("/check-employee-availability", h.CheckEmployeeAvailability)
	r.Get("/check-audience-availability", h.CheckAudienceAvailability)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByEmployee возвращает занятия преподавателя <--->
// [ENG] GetByEmployee returns lessons for the specified employee
func (h *LessonHandler) GetByEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, ok := api.ParseIntParam(w, r, h.Logger, "employee_id")
	if !ok {
		return
	}

	lessons, err := h.manager.GetByEmployee(r.Context(), employeeID)
	if err != nil {
		h.Logger.Error("GetByEmployee failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(lessons),
		len(lessons),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByGroup возвращает занятия группы <--->
// [ENG] GetByGroup returns lessons for the specified group
func (h *LessonHandler) GetByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, ok := api.ParseIntParam(w, r, h.Logger, "group_id")
	if !ok {
		return
	}

	lessons, err := h.manager.GetByGroup(r.Context(), groupID)
	if err != nil {
		h.Logger.Error("GetByGroup failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(lessons),
		len(lessons),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByStudent возвращает индивидуальные занятия студента <--->
// [ENG] GetByStudent returns individual lessons for the specified student
func (h *LessonHandler) GetByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "student_id")
	if !ok {
		return
	}

	lessons, err := h.manager.GetByStudent(r.Context(), studentID)
	if err != nil {
		h.Logger.Error("GetByStudent failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(lessons),
		len(lessons),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetBySubject возвращает занятия по предмету <--->
// [ENG] GetBySubject returns lessons for the specified subject
func (h *LessonHandler) GetBySubject(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := api.ParseIntParam(w, r, h.Logger, "subject_id")
	if !ok {
		return
	}

	lessons, err := h.manager.GetBySubject(r.Context(), subjectID)
	if err != nil {
		h.Logger.Error("GetBySubject failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(lessons),
		len(lessons),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByAudience возвращает занятия в аудитории <--->
// [ENG] GetByAudience returns lessons for the specified audience
func (h *LessonHandler) GetByAudience(w http.ResponseWriter, r *http.Request) {
	audienceID, ok := api.ParseIntParam(w, r, h.Logger, "audience_id")
	if !ok {
		return
	}

	lessons, err := h.manager.GetByAudience(r.Context(), audienceID)
	if err != nil {
		h.Logger.Error("GetByAudience failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(lessons),
		len(lessons),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckEmployeeAvailability проверяет, свободен ли преподаватель в указанное время <--->
// [ENG] CheckEmployeeAvailability checks if the employee is available during the specified time
func (h *LessonHandler) CheckEmployeeAvailability(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(r.URL.Query().Get("employee_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid employee_id")))
		return
	}

	startTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("start_time"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid start_time format")))
		return
	}

	endTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("end_time"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid end_time format")))
		return
	}

	excludeLessonID, _ := strconv.Atoi(r.URL.Query().Get("exclude_lesson_id"))

	isAvailable, err := h.manager.CheckEmployeeAvailability(r.Context(), employeeID, startTime, endTime, excludeLessonID)
	if err != nil {
		h.Logger.Error("CheckEmployeeAvailability failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_available": isAvailable,
		"employee_id":  employeeID,
	})
}

// [RU] CheckAudienceAvailability проверяет, свободна ли аудитория в указанное время <--->
// [ENG] CheckAudienceAvailability checks if the audience is available during the specified time
func (h *LessonHandler) CheckAudienceAvailability(w http.ResponseWriter, r *http.Request) {
	audienceID, err := strconv.Atoi(r.URL.Query().Get("audience_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid audience_id")))
		return
	}

	startTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("start_time"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid start_time format")))
		return
	}

	endTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("end_time"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid end_time format")))
		return
	}

	excludeLessonID, _ := strconv.Atoi(r.URL.Query().Get("exclude_lesson_id"))

	isAvailable, err := h.manager.CheckAudienceAvailability(r.Context(), audienceID, startTime, endTime, excludeLessonID)
	if err != nil {
		h.Logger.Error("CheckAudienceAvailability failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_available": isAvailable,
		"audience_id":  audienceID,
	})
}

// [RU] BulkCreate массово создает занятия <--->
// [ENG] BulkCreate creates multiple lessons
func (h *LessonHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var lessons []*domain.Lesson
	if err := render.DecodeJSON(r.Body, &lessons); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, lesson := range lessons {
		if err := lesson.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), lessons); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

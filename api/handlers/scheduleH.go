package handlers

import (
	"GO_Music/api"
	dto "GO_Music/api/DTO"
	"GO_Music/domain"
	m "GO_Music/engine/managers"
	"errors"
	"net/http"
	"strconv"

	"github.com/SerMoskvin/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// ScheduleHandler обработчик для расписания
type ScheduleHandler struct {
	*api.BaseHandler[int, domain.Schedule, *domain.Schedule,
		dto.ScheduleCreateDTO, dto.ScheduleUpdateDTO, dto.ScheduleResponseDTO]
	manager *m.ScheduleManager
	mapper  *dto.ScheduleMapper
}

// NewScheduleHandler создает новый обработчик расписания
func NewScheduleHandler(
	manager *m.ScheduleManager,
	logger *logger.LevelLogger,
) *ScheduleHandler {
	mapper := dto.NewScheduleMapper()

	return &ScheduleHandler{
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

// Routes возвращает маршруты для расписания
func (h *ScheduleHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Get("/by-lesson/{lesson_id}", h.GetByLesson)
	r.Get("/by-day/{day_week}", h.GetByDay)
	r.Get("/current", h.GetCurrentSchedule)
	r.Get("/check-conflict", h.CheckTimeConflict)
	r.Get("/by-date-range", h.GetByDateRange)
	r.Post("/generate", h.GenerateSchedule)

	return r
}

// [RU] GetByLesson возвращает расписание для конкретного занятия <--->
// [ENG] GetByLesson returns the schedule for a specific lesson
func (h *ScheduleHandler) GetByLesson(w http.ResponseWriter, r *http.Request) {
	lessonID, ok := api.ParseIntParam(w, r, h.Logger, "lesson_id")
	if !ok {
		return
	}

	schedules, err := h.manager.GetByLesson(r.Context(), lessonID)
	if err != nil {
		h.Logger.Error("GetByLesson failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(schedules),
		len(schedules),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByDay возвращает расписание на конкретный день недели <--->
// [ENG] GetByDay returns the schedule for a specific day of the week
func (h *ScheduleHandler) GetByDay(w http.ResponseWriter, r *http.Request) {
	dayWeek := chi.URLParam(r, "day_week")
	if dayWeek == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("day_week is required")))
		return
	}

	schedules, err := h.manager.GetByDay(r.Context(), dayWeek)
	if err != nil {
		h.Logger.Error("GetByDay failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(schedules),
		len(schedules),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetCurrentSchedule возвращает актуальное расписание <--->
// [ENG] GetCurrentSchedule returns the current schedule
func (h *ScheduleHandler) GetCurrentSchedule(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.manager.GetCurrentSchedule(r.Context())
	if err != nil {
		h.Logger.Error("GetCurrentSchedule failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(schedules),
		len(schedules),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckTimeConflict проверяет наличие конфликтов в расписании <--->
// [ENG] CheckTimeConflict checks for conflicts in the schedule
func (h *ScheduleHandler) CheckTimeConflict(w http.ResponseWriter, r *http.Request) {
	dayWeek := r.URL.Query().Get("day_week")
	timeBegin := r.URL.Query().Get("time_begin")
	timeEnd := r.URL.Query().Get("time_end")
	excludeID, _ := strconv.Atoi(r.URL.Query().Get("exclude_id"))

	if dayWeek == "" || timeBegin == "" || timeEnd == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("day_week, time_begin and time_end are required")))
		return
	}

	hasConflict, err := h.manager.CheckTimeConflict(r.Context(), dayWeek, timeBegin, timeEnd, excludeID)
	if err != nil {
		h.Logger.Error("CheckTimeConflict failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"has_conflict": hasConflict,
		"day_week":     dayWeek,
		"time_begin":   timeBegin,
		"time_end":     timeEnd,
	})
}

// [RU] GetByDateRange возвращает расписание в указанном периоде <--->
// [ENG] GetByDateRange returns the schedule in the specified date range
func (h *ScheduleHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date") // Ожидается "DD.MM.YYYY"
	endDateStr := r.URL.Query().Get("end_date")     // Ожидается "DD.MM.YYYY"

	if startDateStr == "" || endDateStr == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("start_date and end_date are required")))
		return
	}

	startDate := domain.ParseDMY(startDateStr)
	endDate := domain.ParseDMY(endDateStr)

	schedules, err := h.manager.GetByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		h.Logger.Error("GetByDateRange failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(schedules),
		len(schedules),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GenerateSchedule генерирует расписание на основе шаблона <--->
// [ENG] GenerateSchedule generates a schedule based on a template
func (h *ScheduleHandler) GenerateSchedule(w http.ResponseWriter, r *http.Request) {
	var template domain.Schedule
	if err := render.DecodeJSON(r.Body, &template); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	untilStr := r.URL.Query().Get("until") // Ожидается "DD.MM.YYYY"
	if untilStr == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("until parameter is required")))
		return
	}

	until := domain.ParseDMY(untilStr)

	if err := template.Validate(); err != nil {
		h.Logger.Error("Validation failed", logger.Error(err))
		render.Render(w, r, api.ErrValidation(err))
		return
	}

	if err := h.manager.GenerateSchedule(r.Context(), &template, until); err != nil {
		h.Logger.Error("GenerateSchedule failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]string{"status": "schedule_generated"})
}

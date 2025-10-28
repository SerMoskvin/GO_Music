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

type StudentAttendanceHandler struct {
	*api.BaseHandler[int, domain.StudentAttendance, *domain.StudentAttendance,
		dto.StudentAttendanceCreateDTO, dto.StudentAttendanceUpdateDTO, dto.StudentAttendanceResponseDTO]
	manager *m.StudentAttendanceManager
	mapper  *dto.StudentAttendanceMapper
}

func NewStudentAttendanceHandler(
	manager *m.StudentAttendanceManager,
	logger *logger.LevelLogger,
) *StudentAttendanceHandler {
	mapper := dto.NewStudentAttendanceMapper()

	return &StudentAttendanceHandler{
		BaseHandler: api.NewBaseHandler(
			manager.BaseManager,
			logger,
			mapper.ToDomain,
			mapper.UpdateDomain,
			mapper.ToResponse,
			nil,
			api.BaseHandlerConfig{
				DefaultPageSize: 50,
				MaxPageSize:     200,
			},
		),
		manager: manager,
		mapper:  mapper,
	}
}

func (h *StudentAttendanceHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/by-student/{student_id}", h.GetByStudent)
	r.Get("/by-lesson/{lesson_id}", h.GetByLesson)
	r.Get("/by-date-range", h.GetByDateRange)
	r.Get("/stats/{student_id}", h.GetStudentAttendanceStats)
	r.Get("/check-duplicate", h.CheckDuplicate)
	r.Post("/bulk-create", h.BulkCreate)

	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	return r
}

// [RU] GetByStudent возвращает записи посещаемости для конкретного студента <--->
// [ENG] GetByStudent returns all student's attendance records
func (h *StudentAttendanceHandler) GetByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "student_id")
	if !ok {
		return
	}

	records, err := h.manager.GetByStudent(r.Context(), studentID)
	if err != nil {
		h.Logger.Error("GetByStudent failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	// Используем форматированные даты для ответа
	response := h.mapper.ToResponseListWithFormattedDate(records)
	api.SendPaginated(w, r, response, len(records), 1, h.Config.DefaultPageSize)
}

// [RU] GetByLesson возвращает записи посещаемости для конкретного занятия <--->
// [ENG] GetByLesson returns all attendance records for a specific lesson
func (h *StudentAttendanceHandler) GetByLesson(w http.ResponseWriter, r *http.Request) {
	lessonID, ok := api.ParseIntParam(w, r, h.Logger, "lesson_id")
	if !ok {
		return
	}

	records, err := h.manager.GetByLesson(r.Context(), lessonID)
	if err != nil {
		h.Logger.Error("GetByLesson failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	response := h.mapper.ToResponseListWithFormattedDate(records)
	api.SendPaginated(w, r, response, len(records), 1, h.Config.DefaultPageSize)
}

// [RU] GetByDateRange возвращает записи за указанный период <--->
// [ENG] GetByDateRange returns attendance records for the specified date range
func (h *StudentAttendanceHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("both start_date and end_date parameters are required")))
		return
	}

	// Преобразуем даты из формата "DD.MM.YYYY" в "YYYY-MM-DD" для БД
	startDateDB := h.formatDateForDB(startDate)
	endDateDB := h.formatDateForDB(endDate)

	if startDateDB == "" || endDateDB == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid date format, expected DD.MM.YYYY")))
		return
	}

	records, err := h.manager.GetByDateRange(r.Context(), startDateDB, endDateDB)
	if err != nil {
		h.Logger.Error("GetByDateRange failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	response := h.mapper.ToResponseListWithFormattedDate(records)
	api.SendPaginated(w, r, response, len(records), 1, h.Config.DefaultPageSize)
}

// [RU] GetStudentAttendanceStats возвращает статистику посещаемости студента <--->
// [ENG] GetStudentAttendanceStats returns the attendance statistics for a student
func (h *StudentAttendanceHandler) GetStudentAttendanceStats(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "student_id")
	if !ok {
		return
	}

	present, absent, err := h.manager.GetStudentAttendanceStats(r.Context(), studentID)
	if err != nil {
		h.Logger.Error("GetStudentAttendanceStats failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	total := present + absent
	attendanceRate := 0.0
	if total > 0 {
		attendanceRate = float64(present) / float64(total) * 100
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"student_id":      studentID,
		"present":         present,
		"absent":          absent,
		"total":           total,
		"attendance_rate": attendanceRate,
	})
}

// [RU] CheckDuplicate проверяет наличие дублирующей записи посещаемости <--->
// [ENG] CheckDuplicate checks for the existence of a duplicate attendance record
func (h *StudentAttendanceHandler) CheckDuplicate(w http.ResponseWriter, r *http.Request) {
	studentID, err := strconv.Atoi(r.URL.Query().Get("student_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid student_id")))
		return
	}

	lessonID, err := strconv.Atoi(r.URL.Query().Get("lesson_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid lesson_id")))
		return
	}

	isDuplicate, err := h.manager.CheckDuplicate(r.Context(), studentID, lessonID)
	if err != nil {
		h.Logger.Error("CheckDuplicate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_duplicate": isDuplicate,
		"student_id":   studentID,
		"lesson_id":    lessonID,
	})
}

// [RU] BulkCreate создает несколько записей посещаемости <--->
// [ENG] BulkCreate creates multiple attendance records
func (h *StudentAttendanceHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var records []*domain.StudentAttendance
	if err := render.DecodeJSON(r.Body, &records); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, record := range records {
		if err := record.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), records); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

// formatDateForDB преобразует дату из "DD.MM.YYYY" в "YYYY-MM-DD"
func (h *StudentAttendanceHandler) formatDateForDB(dateStr string) string {
	if t := domain.ParseDMY(dateStr); !t.IsZero() {
		return t.Format("2006-01-02")
	}
	return ""
}

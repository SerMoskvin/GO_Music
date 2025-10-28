package handlers

import (
	"GO_Music/api"
	dto "GO_Music/api/DTO"
	"GO_Music/domain"
	m "GO_Music/engine/managers"
	"errors"
	"net/http"

	"github.com/SerMoskvin/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type StudentAssessmentHandler struct {
	*api.BaseHandler[int, domain.StudentAssessment, *domain.StudentAssessment,
		dto.StudentAssessmentCreateDTO, dto.StudentAssessmentUpdateDTO, dto.StudentAssessmentResponseDTO]
	manager *m.StudentAssessmentManager
	mapper  *dto.AssessmentMapper
}

func NewStudentAssessmentHandler(
	manager *m.StudentAssessmentManager,
	logger *logger.LevelLogger,
) *StudentAssessmentHandler {
	mapper := dto.NewAssessmentMapper()

	return &StudentAssessmentHandler{
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

func (h *StudentAssessmentHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/by-student/{student_id}", h.GetByStudent)
	r.Get("/by-lesson/{lesson_id}", h.GetByLesson)
	r.Get("/by-task-type/{task_type}", h.GetByTaskType)
	r.Get("/average-grade/{student_id}", h.GetStudentAverageGrade)
	r.Get("/by-date-range", h.GetGradesByDateRange)
	r.Post("/bulk-upsert", h.BulkUpsert)

	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	return r
}

// [RU] GetByStudent возвращает оценки студента <--->
// [ENG] GetByStudent returns student's grades
func (h *StudentAssessmentHandler) GetByStudent(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "student_id")
	if !ok {
		return
	}

	assessments, err := h.manager.GetByStudent(r.Context(), studentID)
	if err != nil {
		h.Logger.Error("GetByStudent failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(assessments),
		len(assessments),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByLesson возвращает оценки за занятие <--->
// [ENG] GetByLesson returns lesson grades
func (h *StudentAssessmentHandler) GetByLesson(w http.ResponseWriter, r *http.Request) {
	lessonID, ok := api.ParseIntParam(w, r, h.Logger, "lesson_id")
	if !ok {
		return
	}

	assessments, err := h.manager.GetByLesson(r.Context(), lessonID)
	if err != nil {
		h.Logger.Error("GetByLesson failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(assessments),
		len(assessments),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByTaskType возвращает оценки по типу задания <--->
// [ENG] GetByTaskType returns grades by task type
func (h *StudentAssessmentHandler) GetByTaskType(w http.ResponseWriter, r *http.Request) {
	taskType := chi.URLParam(r, "task_type")
	if taskType == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("task type is required")))
		return
	}

	assessments, err := h.manager.GetByTaskType(r.Context(), taskType)
	if err != nil {
		h.Logger.Error("GetByTaskType failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(assessments),
		len(assessments),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetStudentAverageGrade возвращает средний балл студента <--->
// [ENG] GetStudentAverageGrade returns student's average grade
func (h *StudentAssessmentHandler) GetStudentAverageGrade(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "student_id")
	if !ok {
		return
	}

	average, err := h.manager.GetStudentAverageGrade(r.Context(), studentID)
	if err != nil {
		h.Logger.Error("GetStudentAverageGrade failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"student_id": studentID,
		"average":    average,
	})
}

// [RU] GetGradesByDateRange возвращает оценки за период <--->
// [ENG] GetGradesByDateRange returns grades for date range
func (h *StudentAssessmentHandler) GetGradesByDateRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("both start_date and end_date are required")))
		return
	}

	// Преобразуем строки DD.MM.YYYY в time.Time для менеджера
	startTime := domain.ParseDMY(startDate)
	endTime := domain.ParseDMY(endDate)

	assessments, err := h.manager.GetGradesByDateRange(r.Context(),
		startTime.Format("2006-01-02"), // Менеджер ожидает формат YYYY-MM-DD
		endTime.Format("2006-01-02"),
	)
	if err != nil {
		h.Logger.Error("GetGradesByDateRange failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(assessments),
		len(assessments),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] BulkUpsert для массового обновления/добавления <--->
// [ENG] BulkUpsert for batch update/insert
func (h *StudentAssessmentHandler) BulkUpsert(w http.ResponseWriter, r *http.Request) {
	var assessments []*domain.StudentAssessment
	if err := render.DecodeJSON(r.Body, &assessments); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, a := range assessments {
		if err := a.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkUpsert(r.Context(), assessments); err != nil {
		h.Logger.Error("BulkUpsert failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

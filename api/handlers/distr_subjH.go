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

type SubjectDistributionHandler struct {
	*api.BaseHandler[int, domain.SubjectDistribution, *domain.SubjectDistribution,
		dto.SubjectDistributionCreateDTO, dto.SubjectDistributionUpdateDTO, dto.SubjectDistributionResponseDTO]
	manager *m.SubjectDistributionManager
	mapper  *dto.SubjectDistributionMapper
}

func NewSubjectDistributionHandler(
	manager *m.SubjectDistributionManager,
	logger *logger.LevelLogger,
) *SubjectDistributionHandler {
	mapper := dto.NewSubjectDistributionMapper()

	return &SubjectDistributionHandler{
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

func (h *SubjectDistributionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Get("/by-employee/{employee_id}", h.GetByEmployee)
	r.Get("/by-subject/{subject_id}", h.GetBySubject)
	r.Get("/by-employee-and-subject", h.GetByEmployeeAndSubject)
	r.Get("/check-exists", h.CheckExists)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByEmployee возвращает распределения по ID сотрудника <--->
// [ENG] GetByEmployee returns distributions by employee ID
func (h *SubjectDistributionHandler) GetByEmployee(w http.ResponseWriter, r *http.Request) {
	employeeID, ok := api.ParseIntParam(w, r, h.Logger, "employee_id")
	if !ok {
		return
	}

	distributions, err := h.manager.GetByEmployee(r.Context(), employeeID)
	if err != nil {
		h.Logger.Error("GetByEmployee failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(distributions),
		len(distributions),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetBySubject возвращает распределения по ID предмета <--->
// [ENG] GetBySubject returns distributions by subject ID
func (h *SubjectDistributionHandler) GetBySubject(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := api.ParseIntParam(w, r, h.Logger, "subject_id")
	if !ok {
		return
	}

	distributions, err := h.manager.GetBySubject(r.Context(), subjectID)
	if err != nil {
		h.Logger.Error("GetBySubject failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(distributions),
		len(distributions),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByEmployeeAndSubject возвращает распределение по ID сотрудника и предмета <--->
// [ENG] GetByEmployeeAndSubject returns distribution by employee ID and subject ID
func (h *SubjectDistributionHandler) GetByEmployeeAndSubject(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(r.URL.Query().Get("employee_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid employee_id")))
		return
	}

	subjectID, err := strconv.Atoi(r.URL.Query().Get("subject_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid subject_id")))
		return
	}

	distribution, err := h.manager.GetByEmployeeAndSubject(r.Context(), employeeID, subjectID)
	if err != nil {
		h.Logger.Error("GetByEmployeeAndSubject failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if distribution == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("distribution not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(distribution))
}

// [RU] CheckExists проверяет существование распределения <--->
// [ENG] CheckExists checks if distribution exists
func (h *SubjectDistributionHandler) CheckExists(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(r.URL.Query().Get("employee_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid employee_id")))
		return
	}

	subjectID, err := strconv.Atoi(r.URL.Query().Get("subject_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid subject_id")))
		return
	}

	exists, err := h.manager.CheckExists(r.Context(), employeeID, subjectID)
	if err != nil {
		h.Logger.Error("CheckExists failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"exists":      exists,
		"employee_id": employeeID,
		"subject_id":  subjectID,
	})
}

// [RU] BulkCreate массово создает распределения <--->
// [ENG] BulkCreate creates multiple distributions
func (h *SubjectDistributionHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var distributions []*domain.SubjectDistribution
	if err := render.DecodeJSON(r.Body, &distributions); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, distr := range distributions {
		if err := distr.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), distributions); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

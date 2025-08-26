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

type EmployeeHandler struct {
	*api.BaseHandler[int, domain.Employee, *domain.Employee,
		dto.EmployeeCreateDTO, dto.EmployeeUpdateDTO, dto.EmployeeResponseDTO]
	manager *m.EmployeeManager
	mapper  *dto.EmployeeMapper
}

func NewEmployeeHandler(
	manager *m.EmployeeManager,
	logger *logger.LevelLogger,
) *EmployeeHandler {
	mapper := dto.NewEmployeeMapper()

	return &EmployeeHandler{
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

func (h *EmployeeHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Get("/by-phone/{phone}", h.GetByPhone)
	r.Get("/by-user/{user_id}", h.GetByUserID)
	r.Get("/by-experience/{min_experience}", h.ListByExperience)
	r.Get("/by-birthday-range", h.ListByBirthdayRange)
	r.Post("/bulk-create", h.BulkCreate)
	r.Get("/check-phone-unique", h.CheckPhoneUnique)

	return r
}

// [RU] GetByPhone возвращает сотрудника по номеру телефона <--->
// [ENG] GetByPhone returns an employee by phone number
func (h *EmployeeHandler) GetByPhone(w http.ResponseWriter, r *http.Request) {
	phone := chi.URLParam(r, "phone")
	if phone == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("phone is required")))
		return
	}

	employee, err := h.manager.GetByPhone(r.Context(), phone)
	if err != nil {
		h.Logger.Error("GetByPhone failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if employee == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("employee not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(employee))
}

// [RU] GetByUserID возвращает сотрудника по ID пользователя <--->
// [ENG] GetByUserID returns an employee by user ID
func (h *EmployeeHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := api.ParseIntParam(w, r, h.Logger, "user_id")
	if !ok {
		return
	}

	employee, err := h.manager.GetByUserID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("GetByUserID failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if employee == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("employee not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(employee))
}

// [RU] ListByExperience возвращает сотрудников с опытом работы не менее указанного <--->
// [ENG] ListByExperience returns employees with at least the specified work experience
func (h *EmployeeHandler) ListByExperience(w http.ResponseWriter, r *http.Request) {
	minExperience, ok := api.ParseIntParam(w, r, h.Logger, "min_experience")
	if !ok {
		return
	}

	employees, err := h.manager.ListByExperience(r.Context(), minExperience)
	if err != nil {
		h.Logger.Error("ListByExperience failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(employees),
		len(employees),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] ListByBirthdayRange возвращает сотрудников с днями рождения в указанном диапазоне <--->
// [ENG] ListByBirthdayRange returns employees with birthdays in the specified range
func (h *EmployeeHandler) ListByBirthdayRange(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("both from and to parameters are required")))
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid from date format")))
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid to date format")))
		return
	}

	employees, err := h.manager.ListByBirthdayRange(r.Context(), from, to)
	if err != nil {
		h.Logger.Error("ListByBirthdayRange failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(employees),
		len(employees),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckPhoneUnique проверяет уникальность номера телефона <--->
// [ENG] CheckPhoneUnique checks the uniqueness of the phone number
func (h *EmployeeHandler) CheckPhoneUnique(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	excludeID, _ := strconv.Atoi(r.URL.Query().Get("exclude_id"))

	if phone == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("phone parameter is required")))
		return
	}

	isUnique, err := h.manager.CheckPhoneUnique(r.Context(), phone, excludeID)
	if err != nil {
		h.Logger.Error("CheckPhoneUnique failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_unique": isUnique,
		"phone":     phone,
	})
}

// [RU] BulkCreate массово создает сотрудников <--->
// [ENG] BulkCreate creates multiple employees
func (h *EmployeeHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var employees []*domain.Employee
	if err := render.DecodeJSON(r.Body, &employees); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, emp := range employees {
		if err := emp.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), employees); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

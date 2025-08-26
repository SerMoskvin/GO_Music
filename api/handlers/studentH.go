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

type StudentHandler struct {
	*api.BaseHandler[int, domain.Student, *domain.Student,
		dto.StudentCreateDTO, dto.StudentUpdateDTO, dto.StudentResponseDTO]
	manager *m.StudentManager
	mapper  *dto.StudentMapper
}

func NewStudentHandler(
	manager *m.StudentManager,
	logger *logger.LevelLogger,
) *StudentHandler {
	mapper := dto.NewStudentMapper()

	return &StudentHandler{
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

func (h *StudentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Get("/by-group/{group_id}", h.GetByGroup)
	r.Get("/by-program/{program_id}", h.GetByProgram)
	r.Get("/search", h.SearchByName)
	r.Get("/by-birthday-range", h.GetByBirthdayRange)
	r.Patch("/{id}/transfer-group", h.TransferToGroup)
	r.Patch("/{id}/change-program", h.ChangeProgram)
	r.Get("/with-account", h.GetWithUserAccount)
	r.Get("/check-phone-unique", h.CheckPhoneNumberUnique)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByGroup возвращает студентов указанной группы <--->
// [ENG] GetByGroup returns students of the specified group
func (h *StudentHandler) GetByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, ok := api.ParseIntParam(w, r, h.Logger, "group_id")
	if !ok {
		return
	}

	students, err := h.manager.GetByGroup(r.Context(), groupID)
	if err != nil {
		h.Logger.Error("GetByGroup failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(students),
		len(students),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByProgram возвращает студентов по программе обучения <--->
// [ENG] GetByProgram returns students by study program
func (h *StudentHandler) GetByProgram(w http.ResponseWriter, r *http.Request) {
	programID, ok := api.ParseIntParam(w, r, h.Logger, "program_id")
	if !ok {
		return
	}

	students, err := h.manager.GetByProgram(r.Context(), programID)
	if err != nil {
		h.Logger.Error("GetByProgram failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(students),
		len(students),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] SearchByName ищет студентов по ФИО <--->
// [ENG] SearchByName searches for students by full name
func (h *StudentHandler) SearchByName(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("search query is required")))
		return
	}

	students, err := h.manager.SearchByName(r.Context(), query)
	if err != nil {
		h.Logger.Error("SearchByName failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(students),
		len(students),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByBirthdayRange возвращает студентов в диапазоне дат рождения <--->
// [ENG] GetByBirthdayRange returns students in the date of birth range
func (h *StudentHandler) GetByBirthdayRange(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("both from and to parameters are required")))
		return
	}

	from := domain.ParseDMY(fromStr)
	to := domain.ParseDMY(toStr)

	students, err := h.manager.GetByBirthdayRange(r.Context(), from, to)
	if err != nil {
		h.Logger.Error("GetByBirthdayRange failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(students),
		len(students),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] TransferToGroup переводит студента в другую группу <--->
// [ENG] TransferToGroup transfers a student to another group
func (h *StudentHandler) TransferToGroup(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "id")
	if !ok {
		return
	}

	var request struct {
		NewGroupID int `json:"new_group_id"`
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	if err := h.manager.TransferToGroup(r.Context(), studentID, request.NewGroupID); err != nil {
		h.Logger.Error("TransferToGroup failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]string{"status": "success"})
}

// [RU] ChangeProgram изменяет программу обучения студента <--->
// [ENG] ChangeProgram changes the student's study program
func (h *StudentHandler) ChangeProgram(w http.ResponseWriter, r *http.Request) {
	studentID, ok := api.ParseIntParam(w, r, h.Logger, "id")
	if !ok {
		return
	}

	var request struct {
		NewProgramID int `json:"new_program_id"`
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	if err := h.manager.ChangeProgram(r.Context(), studentID, request.NewProgramID); err != nil {
		h.Logger.Error("ChangeProgram failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]string{"status": "success"})
}

// [RU] GetWithUserAccount возвращает студентов с привязанными учетными записями <--->
// [ENG] GetWithUserAccount returns students with associated user accounts
func (h *StudentHandler) GetWithUserAccount(w http.ResponseWriter, r *http.Request) {
	students, err := h.manager.GetWithUserAccount(r.Context())
	if err != nil {
		h.Logger.Error("GetWithUserAccount failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(students),
		len(students),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckPhoneNumberUnique проверяет уникальность номера телефона <--->
// [ENG] CheckPhoneNumberUnique checks the uniqueness of the phone number
func (h *StudentHandler) CheckPhoneNumberUnique(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	excludeID, _ := strconv.Atoi(r.URL.Query().Get("exclude_id"))

	if phone == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("phone parameter is required")))
		return
	}

	isUnique, err := h.manager.CheckPhoneNumberUnique(r.Context(), phone, excludeID)
	if err != nil {
		h.Logger.Error("CheckPhoneNumberUnique failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_unique": isUnique,
		"phone":     phone,
	})
}

// [RU] BulkCreate массово создает студентов <--->
// [ENG] BulkCreate creates multiple students
func (h *StudentHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var students []*domain.Student
	if err := render.DecodeJSON(r.Body, &students); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, student := range students {
		if err := student.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), students); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

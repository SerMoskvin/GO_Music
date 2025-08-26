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

type StudyGroupHandler struct {
	*api.BaseHandler[int, domain.StudyGroup, *domain.StudyGroup,
		dto.StudyGroupCreateDTO, dto.StudyGroupUpdateDTO, dto.StudyGroupResponseDTO]
	manager *m.StudyGroupManager
	mapper  *dto.StudyGroupMapper
}

func NewStudyGroupHandler(
	manager *m.StudyGroupManager,
	logger *logger.LevelLogger,
) *StudyGroupHandler {
	mapper := dto.NewStudyGroupMapper()

	return &StudyGroupHandler{
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

func (h *StudyGroupHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Get("/by-program/{program_id}", h.GetByProgram)
	r.Get("/by-name/{name}", h.GetByName)
	r.Get("/by-year/{year}", h.GetByYear)
	r.Get("/check-name-unique", h.CheckNameUnique)
	r.Patch("/{id}/student-count", h.UpdateStudentCount)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByProgram возвращает группы по программе обучения <--->
// [ENG] GetByProgram returns groups by training program
func (h *StudyGroupHandler) GetByProgram(w http.ResponseWriter, r *http.Request) {
	programID, ok := api.ParseIntParam(w, r, h.Logger, "program_id")
	if !ok {
		return
	}

	groups, err := h.manager.GetByProgram(r.Context(), programID)
	if err != nil {
		h.Logger.Error("GetByProgram failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(groups),
		len(groups),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByName возвращает группу по названию <--->
// [ENG] GetByName returns a group by name
func (h *StudyGroupHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("name is required")))
		return
	}

	group, err := h.manager.GetByName(r.Context(), name)
	if err != nil {
		h.Logger.Error("GetByName failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if group == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("group not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(group))
}

// [RU] GetByYear возвращает группы по учебному году <--->
// [ENG] GetByYear returns groups by academic year
func (h *StudyGroupHandler) GetByYear(w http.ResponseWriter, r *http.Request) {
	year, ok := api.ParseIntParam(w, r, h.Logger, "year")
	if !ok {
		return
	}

	groups, err := h.manager.GetByYear(r.Context(), year)
	if err != nil {
		h.Logger.Error("GetByYear failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(groups),
		len(groups),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckNameUnique проверяет уникальность названия группы <--->
// [ENG] CheckNameUnique checks the uniqueness of the group name
func (h *StudyGroupHandler) CheckNameUnique(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	excludeID, _ := strconv.Atoi(r.URL.Query().Get("exclude_id"))

	if name == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("name parameter is required")))
		return
	}

	isUnique, err := h.manager.CheckNameUnique(r.Context(), name, excludeID)
	if err != nil {
		h.Logger.Error("CheckNameUnique failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_unique": isUnique,
		"name":      name,
	})
}

// [RU] UpdateStudentCount обновляет количество студентов в группе <--->
// [ENG] UpdateStudentCount updates the number of students in the group
func (h *StudyGroupHandler) UpdateStudentCount(w http.ResponseWriter, r *http.Request) {
	groupID, ok := api.ParseIntParam(w, r, h.Logger, "id")
	if !ok {
		return
	}

	var request struct {
		NumberOfStudents int `json:"number_of_students" validate:"required"`
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	if err := h.manager.UpdateStudentCount(r.Context(), groupID, request.NumberOfStudents); err != nil {
		h.Logger.Error("UpdateStudentCount failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"status":             "success",
		"group_id":           groupID,
		"number_of_students": request.NumberOfStudents,
	})
}

// [RU] BulkCreate массово создает учебные группы <--->
// [ENG] BulkCreate creates multiple study groups
func (h *StudyGroupHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var groups []*domain.StudyGroup
	if err := render.DecodeJSON(r.Body, &groups); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, group := range groups {
		if err := group.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), groups); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

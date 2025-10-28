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

type SubjectHandler struct {
	*api.BaseHandler[int, domain.Subject, *domain.Subject,
		dto.SubjectCreateDTO, dto.SubjectUpdateDTO, dto.SubjectResponseDTO]
	manager *m.SubjectManager
	mapper  *dto.SubjectMapper
}

func NewSubjectHandler(
	manager *m.SubjectManager,
	logger *logger.LevelLogger,
) *SubjectHandler {
	mapper := dto.NewSubjectMapper()

	return &SubjectHandler{
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

func (h *SubjectHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	r.Get("/by-type/{type}", h.GetByType)
	r.Get("/search-by-name", h.SearchByName)
	r.Get("/search-by-description", h.GetByDescription)
	r.Get("/with-programs/{program_id}", h.GetSubjectsWithPrograms)
	r.Get("/popular", h.GetPopularSubjects)
	r.Get("/check-name-unique", h.CheckNameUnique)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByType возвращает предметы указанного типа <--->
// [ENG] GetByType returns subjects of the specified type
func (h *SubjectHandler) GetByType(w http.ResponseWriter, r *http.Request) {
	subjectType := chi.URLParam(r, "type")
	if subjectType == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("type parameter is required")))
		return
	}

	subjects, err := h.manager.GetByType(r.Context(), subjectType)
	if err != nil {
		h.Logger.Error("GetByType failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(subjects),
		len(subjects),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] SearchByName ищет предметы по названию <--->
// [ENG] SearchByName searches for subjects by name
func (h *SubjectHandler) SearchByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("name parameter is required")))
		return
	}

	subjects, err := h.manager.SearchByName(r.Context(), name)
	if err != nil {
		h.Logger.Error("SearchByName failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(subjects),
		len(subjects),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByDescription ищет предметы по описанию <--->
// [ENG] GetByDescription searches for subjects by description
func (h *SubjectHandler) GetByDescription(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("keyword parameter is required")))
		return
	}

	subjects, err := h.manager.GetByDescription(r.Context(), keyword)
	if err != nil {
		h.Logger.Error("GetByDescription failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(subjects),
		len(subjects),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetSubjectsWithPrograms возвращает предметы с привязанными программами <--->
// [ENG] GetSubjectsWithPrograms returns subjects with associated programs
func (h *SubjectHandler) GetSubjectsWithPrograms(w http.ResponseWriter, r *http.Request) {
	programID, ok := api.ParseIntParam(w, r, h.Logger, "program_id")
	if !ok {
		return
	}

	subjects, err := h.manager.GetSubjectsWithPrograms(r.Context(), programID)
	if err != nil {
		h.Logger.Error("GetSubjectsWithPrograms failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(subjects),
		len(subjects),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetPopularSubjects возвращает самые популярные предметы <--->
// [ENG] GetPopularSubjects returns the most popular subjects
func (h *SubjectHandler) GetPopularSubjects(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	subjects, err := h.manager.GetPopularSubjects(r.Context(), limit)
	if err != nil {
		h.Logger.Error("GetPopularSubjects failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(subjects),
		len(subjects),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckNameUnique проверяет уникальность названия предмета <--->
// [ENG] CheckNameUnique checks the uniqueness of the subject name
func (h *SubjectHandler) CheckNameUnique(w http.ResponseWriter, r *http.Request) {
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

// [RU] BulkCreate массово создает предметы <--->
// [ENG] BulkCreate creates multiple subjects
func (h *SubjectHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var subjects []*domain.Subject
	if err := render.DecodeJSON(r.Body, &subjects); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, subject := range subjects {
		if err := subject.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), subjects); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

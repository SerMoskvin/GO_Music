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

type ProgrammHandler struct {
	*api.BaseHandler[int, domain.Programm, *domain.Programm,
		dto.ProgrammCreateDTO, dto.ProgrammUpdateDTO, dto.ProgrammResponseDTO]
	manager *m.ProgrammManager
	mapper  *dto.ProgrammMapper
}

func NewProgrammHandler(
	manager *m.ProgrammManager,
	logger *logger.LevelLogger,
) *ProgrammHandler {
	mapper := dto.NewProgrammMapper()

	return &ProgrammHandler{
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

func (h *ProgrammHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	r.Get("/by-type/{type}", h.GetByType)
	r.Get("/by-instrument/{instrument}", h.GetByInstrument)
	r.Get("/by-name/{name}", h.GetByName)
	r.Get("/by-duration-range", h.GetByDurationRange)
	r.Get("/by-study-load/{study_load}", h.GetByStudyLoad)
	r.Get("/check-name-unique", h.CheckNameUnique)
	r.Get("/search", h.SearchByDescription)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByType возвращает программы указанного типа <--->
// [ENG] GetByType returns programs of the specified type
func (h *ProgrammHandler) GetByType(w http.ResponseWriter, r *http.Request) {
	programmType := chi.URLParam(r, "type")
	if programmType == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("type is required")))
		return
	}

	programms, err := h.manager.GetByType(r.Context(), programmType)
	if err != nil {
		h.Logger.Error("GetByType failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(programms),
		len(programms),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByInstrument возвращает программы для указанного инструмента <--->
// [ENG] GetByInstrument returns programs for the specified instrument
func (h *ProgrammHandler) GetByInstrument(w http.ResponseWriter, r *http.Request) {
	instrument := chi.URLParam(r, "instrument")
	if instrument == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("instrument is required")))
		return
	}

	programms, err := h.manager.GetByInstrument(r.Context(), instrument)
	if err != nil {
		h.Logger.Error("GetByInstrument failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(programms),
		len(programms),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByName возвращает программу по точному названию <--->
// [ENG] GetByName returns a program by exact name
func (h *ProgrammHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("name is required")))
		return
	}

	programm, err := h.manager.GetByName(r.Context(), name)
	if err != nil {
		h.Logger.Error("GetByName failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if programm == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("programm not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(programm))
}

// [RU] GetByDurationRange возвращает программы в указанном диапазоне длительности <--->
// [ENG] GetByDurationRange returns programs in the specified duration range
func (h *ProgrammHandler) GetByDurationRange(w http.ResponseWriter, r *http.Request) {
	minDurationStr := r.URL.Query().Get("min_duration")
	maxDurationStr := r.URL.Query().Get("max_duration")

	if minDurationStr == "" || maxDurationStr == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("both min_duration and max_duration are required")))
		return
	}

	minDuration, err := strconv.Atoi(minDurationStr)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid min_duration")))
		return
	}

	maxDuration, err := strconv.Atoi(maxDurationStr)
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid max_duration")))
		return
	}

	programms, err := h.manager.GetByDurationRange(r.Context(), minDuration, maxDuration)
	if err != nil {
		h.Logger.Error("GetByDurationRange failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(programms),
		len(programms),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByStudyLoad возвращает программы с указанной учебной нагрузкой <--->
// [ENG] GetByStudyLoad returns programs with the specified study load
func (h *ProgrammHandler) GetByStudyLoad(w http.ResponseWriter, r *http.Request) {
	studyLoad, ok := api.ParseIntParam(w, r, h.Logger, "study_load")
	if !ok {
		return
	}

	programms, err := h.manager.GetByStudyLoad(r.Context(), studyLoad)
	if err != nil {
		h.Logger.Error("GetByStudyLoad failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(programms),
		len(programms),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckNameUnique проверяет уникальность названия программы <--->
// [ENG] CheckNameUnique checks the uniqueness of the program name
func (h *ProgrammHandler) CheckNameUnique(w http.ResponseWriter, r *http.Request) {
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

// [RU] SearchByDescription возвращает программы, содержащие указанный текст в описании <--->
// [ENG] SearchByDescription returns programs containing the specified text in the description
func (h *ProgrammHandler) SearchByDescription(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("q")
	if searchText == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("search query parameter 'q' is required")))
		return
	}

	programms, err := h.manager.SearchByDescription(r.Context(), searchText)
	if err != nil {
		h.Logger.Error("SearchByDescription failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(programms),
		len(programms),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] BulkCreate массово создает программы <--->
// [ENG] BulkCreate creates multiple programs
func (h *ProgrammHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var programms []*domain.Programm
	if err := render.DecodeJSON(r.Body, &programms); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, prog := range programms {
		if err := prog.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), programms); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

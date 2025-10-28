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

type InstrumentHandler struct {
	*api.BaseHandler[int, domain.Instrument, *domain.Instrument,
		dto.InstrumentCreateDTO, dto.InstrumentUpdateDTO, dto.InstrumentResponseDTO]
	manager *m.InstrumentManager
	mapper  *dto.InstrumentMapper
}

func NewInstrumentHandler(
	manager *m.InstrumentManager,
	logger *logger.LevelLogger,
) *InstrumentHandler {
	mapper := dto.NewInstrumentMapper()

	return &InstrumentHandler{
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

func (h *InstrumentHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	r.Get("/by-audience/{audience_id}", h.GetByAudience)
	r.Get("/by-type/{type}", h.GetByType)
	r.Get("/by-name/{name}", h.GetByName)
	r.Get("/check-name-unique", h.CheckNameUnique)
	r.Patch("/{id}/condition", h.UpdateCondition)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByAudience возвращает инструменты в указанной аудитории <--->
// [ENG] GetByAudience returns instruments in the specified audience
func (h *InstrumentHandler) GetByAudience(w http.ResponseWriter, r *http.Request) {
	audienceID, ok := api.ParseIntParam(w, r, h.Logger, "audience_id")
	if !ok {
		return
	}

	instruments, err := h.manager.GetByAudience(r.Context(), audienceID)
	if err != nil {
		h.Logger.Error("GetByAudience failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(instruments),
		len(instruments),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByType возвращает инструменты указанного типа <--->
// [ENG] GetByType returns instruments of the specified type
func (h *InstrumentHandler) GetByType(w http.ResponseWriter, r *http.Request) {
	instrType := chi.URLParam(r, "type")
	if instrType == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("type parameter is required")))
		return
	}

	instruments, err := h.manager.GetByType(r.Context(), instrType)
	if err != nil {
		h.Logger.Error("GetByType failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(instruments),
		len(instruments),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] GetByName возвращает инструмент по точному названию <--->
// [ENG] GetByName returns an instrument by exact name
func (h *InstrumentHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("name parameter is required")))
		return
	}

	instrument, err := h.manager.GetByName(r.Context(), name)
	if err != nil {
		h.Logger.Error("GetByName failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if instrument == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("instrument not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(instrument))
}

// [RU] CheckNameUnique проверяет уникальность названия инструмента <--->
// [ENG] CheckNameUnique checks the uniqueness of the instrument name
func (h *InstrumentHandler) CheckNameUnique(w http.ResponseWriter, r *http.Request) {
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

// [RU] UpdateCondition обновляет состояние инструмента <--->
// [ENG] UpdateCondition updates the condition of the instrument
func (h *InstrumentHandler) UpdateCondition(w http.ResponseWriter, r *http.Request) {
	instrumentID, ok := api.ParseIntParam(w, r, h.Logger, "id")
	if !ok {
		return
	}

	var request struct {
		Condition string `json:"condition" validate:"required,min=1,max=70"`
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	if err := h.manager.UpdateCondition(r.Context(), instrumentID, request.Condition); err != nil {
		h.Logger.Error("UpdateCondition failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]string{"status": "success"})
}

// [RU] BulkCreate массово создает инструменты <--->
// [ENG] BulkCreate creates multiple instruments
func (h *InstrumentHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var instruments []*domain.Instrument
	if err := render.DecodeJSON(r.Body, &instruments); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	for _, instr := range instruments {
		if err := instr.Validate(); err != nil {
			h.Logger.Error("Validation failed", logger.Error(err))
			render.Render(w, r, api.ErrValidation(err))
			return
		}
	}

	if err := h.manager.BulkCreate(r.Context(), instruments); err != nil {
		h.Logger.Error("BulkCreate failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, map[string]string{"status": "success"})
}

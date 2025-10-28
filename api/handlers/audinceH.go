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

type AudienceHandler struct {
	*api.BaseHandler[int, domain.Audience, *domain.Audience,
		dto.AudienceCreateDTO, dto.AudienceUpdateDTO, dto.AudienceResponseDTO]
	manager *m.AudienceManager
	mapper  *dto.AudienceMapper
}

func NewAudienceHandler(
	manager *m.AudienceManager,
	logger *logger.LevelLogger,
) *AudienceHandler {
	mapper := dto.NewAudienceMapper()

	return &AudienceHandler{
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

func (h *AudienceHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.BaseHandler.List)
	r.Post("/", h.BaseHandler.Create)
	r.Get("/{id}", h.BaseHandler.Get)
	r.Put("/{id}", h.BaseHandler.Update)
	r.Patch("/{id}", h.BaseHandler.PartialUpdate)
	r.Delete("/{id}", h.BaseHandler.Delete)

	r.Get("/by-number/{number}", h.GetByNumber)
	r.Get("/by-capacity/{min_capacity}", h.ListByCapacity)
	r.Get("/check-number-unique", h.CheckNumberUnique)

	return r
}

// [RU] GetByNumber возвращает аудиторию по номеру <--->
// [ENG] GetByNumber returns audience by number
func (h *AudienceHandler) GetByNumber(w http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "number")
	if number == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("number is required")))
		return
	}

	audience, err := h.manager.GetByNumber(r.Context(), number)
	if err != nil {
		h.Logger.Error("GetByNumber failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if audience == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("audience not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(audience))
}

// [RU] ListByCapacity возвращает аудитории с вместимостью >= minCapacity <--->
// [ENG] ListByCapacity returns audiences with capacity >= minCapacity
func (h *AudienceHandler) ListByCapacity(w http.ResponseWriter, r *http.Request) {
	minCapacity, ok := api.ParseIntParam(w, r, h.Logger, "min_capacity")
	if !ok {
		return
	}

	audiences, err := h.manager.ListByCapacity(r.Context(), minCapacity)
	if err != nil {
		h.Logger.Error("ListByCapacity failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(audiences),
		len(audiences),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckNumberUnique проверяет уникальность номера аудитории <--->
// [ENG] CheckNumberUnique checks audience number uniqueness
func (h *AudienceHandler) CheckNumberUnique(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	excludeID, _ := strconv.Atoi(r.URL.Query().Get("exclude_id"))

	if number == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("number parameter is required")))
		return
	}

	isUnique, err := h.manager.CheckNumberUnique(r.Context(), number, excludeID)
	if err != nil {
		h.Logger.Error("CheckNumberUnique failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_unique": isUnique,
		"number":    number,
	})
}

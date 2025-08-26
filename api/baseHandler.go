package api

import (
	"net/http"

	"GO_Music/domain"
	"GO_Music/engine"

	"github.com/SerMoskvin/logger"
	"github.com/SerMoskvin/validate"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// BaseHandlerConfig конфигурация для базового обработчика
type BaseHandlerConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

// BaseHandler базовый обработчик для CRUD операций
type BaseHandler[ID comparable, T any, PT interface {
	*T
	domain.Entity[ID]
}, CreateDTO any, UpdateDTO any, ResponseDTO any] struct {
	Manager      *engine.BaseManager[ID, T, PT]
	Logger       *logger.LevelLogger
	ToDomain     func(*CreateDTO) PT     // Маппер из CreateDTO в доменную модель
	UpdateDomain func(PT, *UpdateDTO)    // Функция обновления доменной модели из UpdateDTO
	ToResponse   func(PT) *ResponseDTO   // Маппер из доменной модели в ResponseDTO
	Validate     func(interface{}) error // Функция валидации (может быть кастомной)
	Config       BaseHandlerConfig       // Конфигурация
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler[ID comparable, T any, PT interface {
	*T
	domain.Entity[ID]
}, CreateDTO any, UpdateDTO any, ResponseDTO any](
	manager *engine.BaseManager[ID, T, PT],
	logger *logger.LevelLogger,
	toDomain func(*CreateDTO) PT,
	updateDomain func(PT, *UpdateDTO),
	toResponse func(PT) *ResponseDTO,
	validateFn func(interface{}) error,
	config BaseHandlerConfig,
) *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO] {
	if validateFn == nil {
		validateFn = validate.ValidateStruct
	}

	return &BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]{
		Manager:      manager,
		Logger:       logger,
		ToDomain:     toDomain,
		UpdateDomain: updateDomain,
		ToResponse:   toResponse,
		Validate:     validateFn,
		Config:       config,
	}
}

// Routes возвращает стандартные CRUD маршруты
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Patch("/{id}", h.PartialUpdate)
	r.Delete("/{id}", h.Delete)

	return r
}

// Create обрабатывает создание новой сущности
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) Create(w http.ResponseWriter, r *http.Request) {
	var dto CreateDTO
	if err := render.DecodeJSON(r.Body, &dto); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.Validate(&dto); err != nil {
		h.Logger.Error("Validation failed", logger.Error(err))
		render.Render(w, r, ErrValidation(err))
		return
	}

	entity := h.ToDomain(&dto)
	if err := h.Manager.Create(r.Context(), entity); err != nil {
		h.Logger.Error("Create failed", logger.Error(err))
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, h.ToResponse(entity))
}

// Get обрабатывает получение сущности по ID
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) Get(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Error("Invalid ID", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	entity, err := h.Manager.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("GetByID failed", logger.Error(err), logger.Any("id", id))
		render.Render(w, r, ErrNotFoundOrInternal(err))
		return
	}

	render.JSON(w, r, h.ToResponse(entity))
}

// List обрабатывает получение списка сущностей
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) List(w http.ResponseWriter, r *http.Request) {
	filter, err := h.parseFilter(r)
	if err != nil {
		h.Logger.Error("Failed to parse filter", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	count, err := h.Manager.Count(r.Context(), filter)
	if err != nil {
		h.Logger.Error("Count failed", logger.Error(err))
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	entities, err := h.Manager.List(r.Context(), filter)
	if err != nil {
		h.Logger.Error("List failed", logger.Error(err))
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	response := make([]*ResponseDTO, len(entities))
	for i, entity := range entities {
		response[i] = h.ToResponse(entity)
	}

	currentPage := 1
	if filter.Offset > 0 {
		currentPage = (filter.Offset / filter.Limit) + 1
	}
	totalPages := (count + filter.Limit - 1) / filter.Limit

	render.JSON(w, r, map[string]interface{}{
		"items": response,
		"total": count,
		"pagination": map[string]int{
			"current_page": currentPage,
			"per_page":     filter.Limit,
			"total_pages":  totalPages,
		},
	})
}

// Update обрабатывает полное обновление сущности
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Error("Invalid ID", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var dto UpdateDTO
	if err := render.DecodeJSON(r.Body, &dto); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.Validate(&dto); err != nil {
		h.Logger.Error("Validation failed", logger.Error(err))
		render.Render(w, r, ErrValidation(err))
		return
	}

	entity, err := h.Manager.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("GetByID failed", logger.Error(err), logger.Any("id", id))
		render.Render(w, r, ErrNotFoundOrInternal(err))
		return
	}

	h.UpdateDomain(entity, &dto)
	if err := h.Manager.Update(r.Context(), entity); err != nil {
		h.Logger.Error("Update failed", logger.Error(err))
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, h.ToResponse(entity))
}

// PartialUpdate обрабатывает частичное обновление сущности
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) PartialUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Error("Invalid ID", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var dto UpdateDTO
	if err := render.DecodeJSON(r.Body, &dto); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	entity, err := h.Manager.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("GetByID failed", logger.Error(err), logger.Any("id", id))
		render.Render(w, r, ErrNotFoundOrInternal(err))
		return
	}

	h.UpdateDomain(entity, &dto)
	if err := h.Manager.Update(r.Context(), entity); err != nil {
		h.Logger.Error("Update failed", logger.Error(err))
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, h.ToResponse(entity))
}

// Delete обрабатывает удаление сущности
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(chi.URLParam(r, "id"))
	if err != nil {
		h.Logger.Error("Invalid ID", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := h.Manager.Delete(r.Context(), id); err != nil {
		h.Logger.Error("Delete failed", logger.Error(err), logger.Any("id", id))
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.NoContent(w, r)
}

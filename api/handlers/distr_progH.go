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

type ProgrammDistributionHandler struct {
	*api.BaseHandler[int, domain.ProgrammDistribution, *domain.ProgrammDistribution,
		dto.ProgrammDistributionCreateDTO, dto.ProgrammDistributionUpdateDTO, dto.ProgrammDistributionResponseDTO]
	manager *m.ProgrammDistributionManager
	mapper  *dto.ProgrammDistributionMapper
}

func NewProgrammDistributionHandler(
	manager *m.ProgrammDistributionManager,
	logger *logger.LevelLogger,
) *ProgrammDistributionHandler {
	mapper := dto.NewProgrammDistributionMapper()

	return &ProgrammDistributionHandler{
		BaseHandler: api.NewBaseHandler(
			manager.BaseManager,
			logger,
			mapper.ToDomain,
			mapper.UpdateDomain,
			mapper.ToResponse,
			nil,
			api.BaseHandlerConfig{
				DefaultPageSize: 50,
				MaxPageSize:     200,
			},
		),
		manager: manager,
		mapper:  mapper,
	}
}

func (h *ProgrammDistributionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Get("/by-programm/{programm_id}", h.GetByProgramm)
	r.Get("/by-subject/{subject_id}", h.GetBySubject)
	r.Get("/check-exists", h.CheckExists)
	r.Get("/by-programm-and-subject", h.GetByProgrammAndSubject)
	r.Post("/bulk-create", h.BulkCreate)

	return r
}

// [RU] GetByProgramm возвращает распределения по ID программы <--->
// [ENG] GetByProgramm returns distributions by program ID
func (h *ProgrammDistributionHandler) GetByProgramm(w http.ResponseWriter, r *http.Request) {
	programmID, ok := api.ParseIntParam(w, r, h.Logger, "programm_id")
	if !ok {
		return
	}

	distributions, err := h.manager.GetByProgramm(r.Context(), programmID)
	if err != nil {
		h.Logger.Error("GetByProgramm failed", logger.Error(err))
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
func (h *ProgrammDistributionHandler) GetBySubject(w http.ResponseWriter, r *http.Request) {
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

// [RU] CheckExists проверяет существование распределения программы <--->
// [ENG] CheckExists checks if program distribution exists
func (h *ProgrammDistributionHandler) CheckExists(w http.ResponseWriter, r *http.Request) {
	programmID, err := strconv.Atoi(r.URL.Query().Get("programm_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid programm_id")))
		return
	}

	subjectID, err := strconv.Atoi(r.URL.Query().Get("subject_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid subject_id")))
		return
	}

	exists, err := h.manager.CheckExists(r.Context(), programmID, subjectID)
	if err != nil {
		h.Logger.Error("CheckExists failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"exists":      exists,
		"programm_id": programmID,
		"subject_id":  subjectID,
	})
}

// [RU] GetByProgrammAndSubject возвращает распределение программы по ID программы и предмета <--->
// [ENG] GetByProgrammAndSubject returns program distribution by program ID and subject ID
func (h *ProgrammDistributionHandler) GetByProgrammAndSubject(w http.ResponseWriter, r *http.Request) {
	programmID, err := strconv.Atoi(r.URL.Query().Get("programm_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid programm_id")))
		return
	}

	subjectID, err := strconv.Atoi(r.URL.Query().Get("subject_id"))
	if err != nil {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid subject_id")))
		return
	}

	distribution, err := h.manager.GetByProgrammAndSubject(r.Context(), programmID, subjectID)
	if err != nil {
		h.Logger.Error("GetByProgrammAndSubject failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if distribution == nil {
		render.Render(w, r, api.ErrNotFoundOrInternal(errors.New("distribution not found")))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(distribution))
}

// [RU] BulkCreate создает несколько распределений программы <--->
// [ENG] BulkCreate creates multiple program distributions
func (h *ProgrammDistributionHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	var distributions []*domain.ProgrammDistribution
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

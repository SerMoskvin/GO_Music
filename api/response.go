package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/SerMoskvin/logger"
	"github.com/SerMoskvin/validate"
	"github.com/go-chi/render"
)

// [RU] ErrResponse стандартная структура для ошибок API <--->
// [ENG] ErrResponse standard API error structure
type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string            `json:"status"`
	AppCode    int64             `json:"code,omitempty"`
	ErrorText  string            `json:"error,omitempty"`
	Validation map[string]string `json:"validation,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// [RU] ErrInvalidRequest создает ответ для невалидных запросов (400) <--->
// [ENG] ErrInvalidRequest creates response for invalid requests (400)
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

// [RU] ErrValidation создает ответ для ошибок валидации (422) <--->
// [ENG] ErrValidation creates response for validation errors (422)
func ErrValidation(err error) render.Renderer {
	if verr, ok := err.(validate.ValidationErrors); ok {
		validation := make(map[string]string)
		for field, details := range verr {
			validation[field] = details.Message
		}
		return &ErrResponse{
			Err:            err,
			HTTPStatusCode: 422,
			StatusText:     "Validation failed",
			Validation:     validation,
		}
	}
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Validation failed",
		ErrorText:      err.Error(),
	}
}

// [RU] ErrNotFoundOrInternal создает ответ для отсутствующих ресурсов (404) или внутренних ошибок (500) <--->
// [ENG] ErrNotFoundOrInternal creates response for not found (404) or internal errors (500)
func ErrNotFoundOrInternal(err error) render.Renderer {
	if errors.Is(err, sql.ErrNoRows) {
		return &ErrResponse{
			Err:            err,
			HTTPStatusCode: 404,
			StatusText:     "Resource not found",
		}
	}
	return ErrInternalServer(err)
}

// [RU] ErrInternalServer создает ответ для внутренних ошибок сервера (500) <--->
// [ENG] ErrInternalServer creates response for internal server errors (500)
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error",
		ErrorText:      err.Error(),
	}
}

// [RU] SendSuccess отправляет успешный JSON ответ (200) <--->
// [ENG] SendSuccess sends successful JSON response (200)
func SendSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

// [RU] SendCreated отправляет ответ о создании ресурса (201) <--->
// [ENG] SendCreated sends resource creation response (201)
func SendCreated(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

// [RU] SendPaginated отправляет пагинированный список <--->
// [ENG] SendPaginated sends paginated list
func SendPaginated(
	w http.ResponseWriter,
	r *http.Request,
	items interface{},
	total int,
	page int,
	pageSize int,
) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"items": items,
		"total": total,
		"pagination": map[string]int{
			"page":     page,
			"per_page": pageSize,
			"pages":    (total + pageSize - 1) / pageSize,
		},
	})
}
func ProcessBody(w http.ResponseWriter, r *http.Request, log *logger.LevelLogger, target interface{}) bool {
	if err := render.DecodeJSON(r.Body, &target); err != nil {
		log.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return false
	}
	return true
}

func Validate(w http.ResponseWriter, r *http.Request, log *logger.LevelLogger, data interface{}, validateFunc func() error) bool {
	if err := validateFunc(); err != nil {
		log.Error("Validation failed", logger.Error(err))
		render.Render(w, r, ErrValidation(err))
		return false
	}
	return true
}

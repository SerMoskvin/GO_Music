package api

import (
	"GO_Music/db"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/SerMoskvin/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// [RU] parseID преобразует строковый ID в нужный тип <--->
// [ENG] parseID converts string ID to required type
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) parseID(idStr string) (ID, error) {
	var zero ID
	t := reflect.TypeOf(zero)

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		idVal, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return zero, errors.New("invalid ID format")
		}
		return reflect.ValueOf(idVal).Convert(t).Interface().(ID), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		idVal, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return zero, errors.New("invalid ID format")
		}
		return reflect.ValueOf(idVal).Convert(t).Interface().(ID), nil

	case reflect.String:
		return reflect.ValueOf(idStr).Convert(t).Interface().(ID), nil

	default:
		return zero, errors.New("unsupported ID type")
	}
}

// [RU] parseFilter парсит параметры запроса в объект фильтра <--->
// [ENG] parseFilter parses request parameters into filter object
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) parseFilter(r *http.Request) (db.Filter, error) {
	query := r.URL.Query()
	filter := db.Filter{}

	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 || pageSize > h.Config.MaxPageSize {
			return filter, errors.New("invalid page_size parameter")
		}
		filter.Limit = pageSize
	} else {
		filter.Limit = h.Config.DefaultPageSize
	}

	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			return filter, errors.New("invalid page parameter")
		}
		filter.Offset = (page - 1) * filter.Limit
	}

	if sort := query.Get("sort"); sort != "" {
		filter.OrderBy = sort
	}

	for key, values := range query {
		if len(values) == 0 || values[0] == "" {
			continue
		}

		switch {
		case strings.HasPrefix(key, "filter["):
			field := strings.TrimPrefix(key, "filter[")
			field = strings.TrimSuffix(field, "]")
			filter.Conditions = append(filter.Conditions, db.Condition{
				Field:    field,
				Operator: "=",
				Value:    values[0],
			})

		case strings.HasPrefix(key, "range["):
			field := strings.TrimPrefix(key, "range[")
			field = strings.TrimSuffix(field, "]")
			parts := strings.Split(field, "][")
			if len(parts) == 2 {
				switch parts[1] {
				case "from":
					if t, err := time.Parse(time.RFC3339, values[0]); err == nil {
						filter.Conditions = append(filter.Conditions, db.Condition{
							Field:    parts[0],
							Operator: ">=",
							Value:    t,
						})
					}
				case "to":
					if t, err := time.Parse(time.RFC3339, values[0]); err == nil {
						filter.Conditions = append(filter.Conditions, db.Condition{
							Field:    parts[0],
							Operator: "<=",
							Value:    t,
						})
					}
				}
			}

		case key == "search":
			filter.Search = values[0]
		}
	}

	return filter, nil
}

// [RU] parseIDFromRequest универсальный парсер ID из URL параметров <--->
// [ENG] parseIDFromRequest universal ID parser from URL parameters
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) parseIDFromRequest(r *http.Request, paramName string) (ID, bool) {
	idStr := chi.URLParam(r, paramName)
	id, err := h.parseID(idStr)
	if err != nil {
		h.Logger.Error("Invalid ID", logger.Error(err), logger.String("param", paramName), logger.String("value", idStr))
		return id, false
	}
	return id, true
}

// [RU] processListResult унифицированная обработка списковых результатов
// [ENG] processListResult unified list results processing
func (h *BaseHandler[ID, T, PT, CreateDTO, UpdateDTO, ResponseDTO]) processListResult(w http.ResponseWriter, r *http.Request, items []PT, total int) {
	response := make([]*ResponseDTO, len(items))
	for i, item := range items {
		response[i] = h.ToResponse(item)
	}

	render.JSON(w, r, map[string]interface{}{
		"items": response,
		"total": total,
	})
}

// [RU] ParseIntParam парсит целочисленный параметр из URL <--->
// [ENG] ParseIntParam parse INT parameter from URL
func ParseIntParam(w http.ResponseWriter, r *http.Request, log *logger.LevelLogger, paramName string) (int, bool) {
	param := chi.URLParam(r, paramName)
	value, err := strconv.Atoi(param)
	if err != nil {
		log.Error("Invalid param",
			logger.Error(err),
			logger.String("param", paramName),
			logger.String("value", param),
		)
		render.Render(w, r, ErrInvalidRequest(err))
		return 0, false
	}
	return value, true
}

// [RU] ParseStringParam проверяет обязательный строковый параметр <--->
// [ENG] ParseStringParam check required string parametr
func ParseStringParam(w http.ResponseWriter, r *http.Request, log *logger.LevelLogger, paramName string) (string, bool) {
	param := chi.URLParam(r, paramName)
	if param == "" {
		err := errors.New(paramName + " is required")
		log.Error("Missing param", logger.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
		return "", false
	}
	return param, true
}

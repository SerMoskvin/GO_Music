package api

import (
	"github.com/SerMoskvin/access"
	"github.com/go-chi/chi/v5"
)

// SetupEntity универсальная функция для настройки роутов любой сущности
func SetupEntity[T interface{ Routes() chi.Router }](router chi.Router, auth *access.Authenticator, handler T, path string) {
	router.Route(path, func(r chi.Router) {
		r.Use(auth.CheckPermissions)
		r.Mount("/", handler.Routes())
	})
}

// SetupAll универсальная функция для настройки всех роутов
func SetupAll(router chi.Router, auth *access.Authenticator, handlers map[string]interface{ Routes() chi.Router }) {
	for path, handler := range handlers {
		SetupEntity(router, auth, handler, "/"+path)
	}
}

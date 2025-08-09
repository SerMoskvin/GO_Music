package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

type BaseHandler struct {
	SessionStore sessions.Store
}

// Получить сессию с обработкой ошибки
func (h *BaseHandler) GetSession(r *http.Request, name string) (*sessions.Session, error) {
	return h.SessionStore.Get(r, name)
}

// Универсальный метод для отправки JSON-ответа
func (h *BaseHandler) WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Универсальный метод для отправки ошибки в JSON
func (h *BaseHandler) WriteError(w http.ResponseWriter, status int, message string) {
	h.WriteJSON(w, status, map[string]string{"error": message})
}

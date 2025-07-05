package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"

	"GO_Music/engine"

	"github.com/SerMoskvin/access"
)

type UserHandler struct {
	*BaseHandler

	UserManager    *engine.UserManager
	UserRepository engine.UserRepository
}

func NewUserHandler(userManager *engine.UserManager, userRepo engine.UserRepository, sessionStore sessions.Store) *UserHandler {
	return &UserHandler{
		BaseHandler:    &BaseHandler{SessionStore: sessionStore},
		UserManager:    userManager,
		UserRepository: userRepo,
	}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token       string `json:"token,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Error       string `json:"error,omitempty"`
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.WriteError(w, http.StatusMethodNotAllowed, "только POST")
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "неверный JSON")
		return
	}

	token, err := h.UserManager.Login(req.Login, req.Password)
	if err != nil {
		h.WriteJSON(w, http.StatusUnauthorized, loginResponse{Error: err.Error()})
		return
	}

	user, err := h.UserRepository.GetByLogin(req.Login)
	if err != nil || user == nil {
		h.WriteJSON(w, http.StatusInternalServerError, loginResponse{Error: "не удалось получить данные пользователя"})
		return
	}

	session, err := h.GetSession(r, "session-name")
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, "не удалось получить сессию")
		return
	}

	session.Values["jwt"] = token
	if err := session.Save(r, w); err != nil {
		h.WriteError(w, http.StatusInternalServerError, "не удалось сохранить сессию")
		return
	}

	redirectURL := access.RedirectURL(user.Role)

	resp := loginResponse{
		Token:       token,
		RedirectURL: redirectURL,
	}
	h.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.GetSession(r, "session-name")
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, "не удалось получить сессию")
		return
	}

	delete(session.Values, "jwt")
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		h.WriteError(w, http.StatusInternalServerError, "не удалось сохранить сессию")
		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"message": "выход выполнен"})
}

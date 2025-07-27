package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"GO_Music/domain"

)


// CreateUserHandler — POST /users
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.WriteError(w, http.StatusMethodNotAllowed, "только POST")
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.WriteError(w, http.StatusBadRequest, "неверный JSON")
		return
	}

	err := h.UserManager.Create(&user)
	if err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusCreated, user)
}

// GetUserHandler — GET /users/{id}
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.WriteError(w, http.StatusBadRequest, "неверный ID")
		return
	}

	user, err := h.UserManager.GetByID(id)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		h.WriteError(w, http.StatusNotFound, "пользователь не найден")
		return
	}

	h.WriteJSON(w, http.StatusOK, user)
}

// UpdateUserHandler — PUT /users/{id}
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.WriteError(w, http.StatusMethodNotAllowed, "только PUT")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.WriteError(w, http.StatusBadRequest, "неверный ID")
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.WriteError(w, http.StatusBadRequest, "неверный JSON")
		return
	}

	user.UserID = id

	err = h.UserManager.Update(&user)
	if err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, user)
}

// DeleteUserHandler — DELETE /users/{id}
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.WriteError(w, http.StatusMethodNotAllowed, "только DELETE")
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		h.WriteError(w, http.StatusBadRequest, "неверный ID")
		return
	}

	err = h.UserManager.Delete(id)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"message": "пользователь удалён"})
}

// GetUsersByIDsHandler — GET /users?ids=1,2,3
func (h *UserHandler) GetUsersByIDsHandler(w http.ResponseWriter, r *http.Request) {
	idsStr := r.URL.Query().Get("ids")
	if idsStr == "" {
		h.WriteError(w, http.StatusBadRequest, "параметр ids обязателен")
		return
	}

	idStrs := strings.Split(idsStr, ",")
	var ids []int
	for _, s := range idStrs {
		id, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || id <= 0 {
			h.WriteError(w, http.StatusBadRequest, "неверный ID в списке")
			return
		}
		ids = append(ids, id)
	}

	users, err := h.UserManager.GetByIDs(ids)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, users)
}

package handlers

import (
	"GO_Music/api"
	dto "GO_Music/api/DTO"
	"GO_Music/domain"
	"GO_Music/engine"
	m "GO_Music/engine/managers"
	"errors"
	"net/http"
	"strconv"

	"github.com/SerMoskvin/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	*api.BaseHandler[int, domain.User, *domain.User,
		dto.UserCreateDTO, dto.UserUpdateDTO, dto.UserResponseDTO]
	manager *m.UserManager
	mapper  *dto.UserMapper
}

func NewUserHandler(
	manager *m.UserManager,
	logger *logger.LevelLogger,
) *UserHandler {
	mapper := dto.NewUserMapper()

	return &UserHandler{
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

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Mount("/", h.BaseHandler.Routes())

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Get("/current", h.GetCurrentUser)
	r.Put("/change-password", h.ChangePassword)
	r.Get("/by-role/{role}", h.GetByRole)
	r.Get("/search", h.SearchByNames)
	r.Get("/check-login-unique", h.CheckLoginUnique)
	r.Post("/{user_id}/image", h.UploadImage) // Простая загрузка
	r.Get("/{user_id}/image", h.GetImage)     // Получение изображения

	return r
}

// [RU] Register создает нового пользователя <--->
// [ENG] Register creates a new user
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var createDTO dto.UserCreateDTO
	if err := render.DecodeJSON(r.Body, &createDTO); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	user := h.mapper.ToDomain(&createDTO)
	if err := h.manager.Register(r.Context(), user); err != nil {
		h.Logger.Error("Register failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendCreated(w, r, h.mapper.ToResponse(user))
}

// [RU] Login выполняет аутентификацию пользователя <--->
// [ENG] Login authenticates the user
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.UserLoginDTO
	if err := render.DecodeJSON(r.Body, &loginDTO); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	token, err := h.manager.Login(r.Context(), loginDTO.Login, loginDTO.Password)
	if err != nil {
		h.Logger.Error("Login failed", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(errors.New("invalid credentials")))
		return
	}

	api.SendSuccess(w, r, map[string]string{
		"token":   token,
		"message": "Login successful",
	})
}

// [RU] GetCurrentUser возвращает данные текущего пользователя <--->
// [ENG] GetCurrentUser returns current user data
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.manager.GetCurrentUser(r.Context())
	if err != nil {
		h.Logger.Error("GetCurrentUser failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, h.mapper.ToResponse(user))
}

// [RU] ChangePassword изменяет пароль пользователя <--->
// [ENG] ChangePassword changes user password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := h.manager.GetCurrentUser(r.Context())
	if err != nil {
		h.Logger.Error("GetCurrentUser failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	var passwordDTO dto.UserChangePasswordDTO
	if err := render.DecodeJSON(r.Body, &passwordDTO); err != nil {
		h.Logger.Error("Failed to decode request body", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	if err := h.manager.ChangePassword(r.Context(), user.UserID, passwordDTO.OldPassword, passwordDTO.NewPassword); err != nil {
		h.Logger.Error("ChangePassword failed", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}

	api.SendSuccess(w, r, map[string]string{"message": "Password changed successfully"})
}

// [RU] GetByRole возвращает пользователей по роли <--->
// [ENG] GetByRole returns users by role
func (h *UserHandler) GetByRole(w http.ResponseWriter, r *http.Request) {
	role := chi.URLParam(r, "role")
	if role == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("role is required")))
		return
	}

	users, err := h.manager.GetByRole(r.Context(), role)
	if err != nil {
		h.Logger.Error("GetByRole failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(users),
		len(users),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] SearchByNames ищет пользователей по ФИО <--->
// [ENG] SearchByNames searches users by full name
func (h *UserHandler) SearchByNames(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("search query is required")))
		return
	}

	users, err := h.manager.SearchByNames(r.Context(), query)
	if err != nil {
		h.Logger.Error("SearchByNames failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendPaginated(w, r,
		h.mapper.ToResponseList(users),
		len(users),
		1,
		h.Config.DefaultPageSize,
	)
}

// [RU] CheckLoginUnique проверяет уникальность логина <--->
// [ENG] CheckLoginUnique checks login uniqueness
func (h *UserHandler) CheckLoginUnique(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	excludeID, _ := strconv.Atoi(r.URL.Query().Get("exclude_id"))

	if login == "" {
		render.Render(w, r, api.ErrInvalidRequest(errors.New("login parameter is required")))
		return
	}

	isUnique, err := h.manager.CheckLoginUnique(r.Context(), login, excludeID)
	if err != nil {
		h.Logger.Error("CheckLoginUnique failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]interface{}{
		"is_unique": isUnique,
		"login":     login,
	})
}

// [RU] UploadImage загружает изображение для пользователя <--->
// [ENG] UploadImage uploads image for user
func (h *UserHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := api.ParseIntParam(w, r, h.Logger, "user_id")
	if !ok {
		return
	}

	// Ограничение размера файла (5MB)
	r.ParseMultipartForm(5 << 20)
	file, header, err := r.FormFile("image")
	if err != nil {
		h.Logger.Error("Failed to get image file", logger.Error(err))
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	defer file.Close()

	imageData := make([]byte, header.Size)
	if _, err := file.Read(imageData); err != nil {
		h.Logger.Error("Failed to read image data", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	user, err := h.manager.GetByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("GetByID failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	user.Image = imageData
	if err := h.manager.Update(r.Context(), user); err != nil {
		h.Logger.Error("Update failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	api.SendSuccess(w, r, map[string]string{"message": "Image uploaded successfully"})
}

// [RU] DownloadImage скачивает изображение пользователя <--->
// [ENG] DownloadImage downloads user image
func (h *UserHandler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := api.ParseIntParam(w, r, h.Logger, "user_id")
	if !ok {
		return
	}

	user, err := h.manager.GetByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("GetByID failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if len(user.Image) == 0 {
		user.Image = engine.DefaultImage
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename=user_"+strconv.Itoa(userID)+".jpg")
	w.Write(user.Image)
}

// [RU] GetImage возвращает изображение пользователя <--->
// [ENG] GetImage returns user image
func (h *UserHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	userID, ok := api.ParseIntParam(w, r, h.Logger, "user_id")
	if !ok {
		return
	}

	user, err := h.manager.GetByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("GetByID failed", logger.Error(err))
		render.Render(w, r, api.ErrInternalServer(err))
		return
	}

	if len(user.Image) == 0 {
		user.Image = engine.DefaultImage
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(user.Image)
}

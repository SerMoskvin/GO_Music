package app

import (
	"log"
	"net/http"

	"GO_Music/engine"
	"GO_Music/handler"
	"GO_Music/repository"
	"GO_Music/routes"
)

// Run инициализирует все компоненты и запускает HTTP-сервер
func Run(addr string) error {
	userRepo, err := repository.NewUserRepository(/* параметры подключения */)
	if err != nil {
		return err
	}


	userManager := engine.NewUserManager(userRepo, /* аутентификатор */)

	// Создание обработчика
	jwtSecret := "секрет_для_jwt" // лучше брать из конфига/окружения
	h := handler.NewHandler(userManager, jwtSecret)
	h.UserRepository = userRepo

	// Регистрируем маршруты
	routes.RegisterRoutes(h)

	// Отдача статики
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Сервер запущен на %s", addr)
	return http.ListenAndServe(addr, nil)
}
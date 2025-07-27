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
	log := logger.NewLevel(logger.DefaultLevelConfig())
}
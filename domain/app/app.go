package app

import (
	"passbase-hometest/domain/database"
)

type App struct {
	Router *Router
}

func NewApp(repository database.Repository) *App {
	return &App{
		Router:     NewRouter(repository),
	}
}

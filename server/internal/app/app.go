package app

import (
	"database/sql"
	"net/http"

	"github.com/ijry/sanjie/internal/handler"
	"github.com/ijry/sanjie/internal/httpx"
)

type App struct {
	DB     *sql.DB
	Router http.Handler
}

func New(db *sql.DB) *App {
	return &App{
		DB:     db,
		Router: httpx.CORS(httpx.DemoAuth(handler.Routes(db))),
	}
}

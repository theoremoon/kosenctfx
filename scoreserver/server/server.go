package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/theoremoon/kosenctfx/scoreserver/service"
)

type Server interface {
	Start(addr string) error
}
type server struct {
	app service.App
}

func New(app service.App) Server {
	return &server{
		app: app,
	}
}

func (s *server) Start(addr string) error {
	e := echo.New()
	e.Use(middleware.Logger())
	return e.Start(addr)
}

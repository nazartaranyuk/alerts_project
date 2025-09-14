package server

import (
	"context"
	"fmt"
	"nazartaraniuk/alertsProject/internal/adapter/handler"
	"nazartaraniuk/alertsProject/internal/adapter/midl"
	"nazartaraniuk/alertsProject/internal/adapter/ws"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	cfg    *config.Config
	server *echo.Echo
}

func NewServer(cfg *config.Config, s usecase.GetAlarmInfoService) (*Server, error) {
	hub := ws.NewHub()
	server := echo.New()

	server.GET("/health", handler.Health())

	server.GET("/location", handler.Handler(hub))

	server.GET("/alerts", handler.GetAlarms(s))

	server.GET("/swagger", echoSwagger.WrapHandler)

	midl.AddTestAuthMiddleWare(server, cfg.Server.AdminUsername, cfg.Server.AdminPassword)

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "Not found")
	})

	return &Server{cfg: cfg, server: server}, nil
}

func (a *Server) Run() error {
	errCh := make(chan error, 1)
	port := fmt.Sprintf(":%d", a.cfg.Server.Port)
	go func() { errCh <- a.server.Start(port) }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return a.server.Shutdown(ctx)
	}
}

package server

import (
	"context"
	"fmt"
	"nazartaraniuk/alertsProject/internal/adapter/handler"
	"nazartaraniuk/alertsProject/internal/adapter/midl"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	cfg    config.Config
	server *echo.Echo
}

func NewServer(cfg config.Config, alarmsService *usecase.GetAlarmInfoService, userService *usecase.UserService) (*Server, error) {
	server := echo.New()
	// Middlewares
	server.Use(middleware.Logger())
	midl.AddJWTMiddleware(server, []byte(cfg.Server.JWTSecret))
	midl.AddPrometheusMiddleware(server)

	// Routes
	server.POST("/api/register", handler.RegistrationHandler(*userService))
	server.POST("/api/login", handler.LoginHandler(cfg, *userService))
	server.POST("/api/send-from-telegram", handler.SendFromBotHandler())

	server.GET("/api/v1/health", handler.Health())
	server.GET("/api/v1/alerts", handler.GetAlarms(alarmsService))
	server.GET("/api/swagger", echoSwagger.WrapHandler)
	server.GET("/api/metrics", echo.WrapHandler(promhttp.Handler()))
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "Wrong endpoint. Please use /api")
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

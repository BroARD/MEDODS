package server

import (
	"Medods/config"
	"Medods/pkg/logging"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	db     *gorm.DB
	logger logging.Logger
}

func NewServer(cfg *config.Config, db *gorm.DB, logger logging.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, logger: logger, db: db}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}
	go func() {
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatal("Error starting server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return s.echo.Server.Shutdown(ctx)
}


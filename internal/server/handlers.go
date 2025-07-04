package server

import (
	authHttp "Medods/internal/auth/delivery/http"
	authRepo "Medods/internal/auth/repository"
	authUseCase "Medods/internal/auth/usecase"
	"Medods/internal/middleware"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	
)

func (s *Server) MapHandlers(e *echo.Echo) error{
	authRepository := authRepo.NewAuthRepository(s.db)
	authUseCase := authUseCase.NewAuthUseCase(authRepository, s.logger)
	authHandlers := authHttp.NewAuthHandlers(authUseCase, s.logger)

	mw := middleware.NewMiddlewareManager(authUseCase, s.cfg, s.logger)

	v1 := e.Group("/api")

	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(echoMid.CORS())

	authGroup := v1.Group("")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	routes := e.Routes()

	for _, route := range routes {
		s.logger.Infof("Method: %s, Path: %s\n", route.Method, route.Path)
	}
	return nil
}
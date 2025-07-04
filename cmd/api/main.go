package main

import (
	"Medods/config"
	"Medods/internal/server"
	"Medods/pkg/db"
	"Medods/pkg/logging"
	"log"

	_ "Medods/docs"
	"github.com/sirupsen/logrus"
)

// @title JWT API
// @version 1.0

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.GetConfig()
	logger, err := logging.NewLogger("logs", logrus.TraceLevel)
	if err != nil {
		log.Fatal(err)
	}

	psqlDB, err := db.InitDB()
	if err != nil {
		logger.Info("Could not start a DB")
		logger.Fatal(err)
	}

	s := server.NewServer(cfg, psqlDB, *logger)
	if err = s.Run(); err != nil {
		logger.Fatal(err)
	}
}

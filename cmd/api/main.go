package main

import (
	"Medods/config"
	"Medods/internal/server"
	"Medods/pkg/db"
	"Medods/pkg/logging"
	"log"

	"github.com/sirupsen/logrus"
)

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

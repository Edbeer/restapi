package main

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/Edbeer/restapi/pkg/db/postgres"
	server "github.com/Edbeer/restapi/internal/transport/rest"
)

func main() {
	cfg := config.GetConfig()
	logger := logger.NewApiLogger(cfg)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		logger.Fatalf("Postgresql init: %s", err)
	} else {
		logger.Infof("Postgres connected, Status: %#v", psqlDB.Stat())
	}
	defer psqlDB.Close()


	logger.InitLogger()
	logger.Info("Starting auth server")
	s := server.NewServer(cfg, psqlDB, logger)
	if err := s.Run(); err != nil {
		logger.Fatal(err)
	}
}
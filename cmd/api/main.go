package main

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/server"
	"github.com/Edbeer/restapi/pkg/db/postgres"
	"github.com/Edbeer/restapi/pkg/db/redis"
	"github.com/Edbeer/restapi/pkg/logger"
)

func main() {
	cfg := config.GetConfig()
	logger := logger.NewApiLogger(cfg)

	// postgresql
	psqlClient, err := postgres.NewPsqlClient(cfg)
	if err != nil {
		logger.Fatalf("Postgresql init: %s", err)
	} else {
		logger.Infof("Postgres connected, Status: %#v", psqlClient.Stat())
	}
	defer psqlClient.Close()

	// redis
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	logger.Info("Redis connetcted")

	logger.InitLogger()
	logger.Info("Starting auth server")
	s := server.NewServer(cfg, psqlClient, logger)
	if err := s.Run(); err != nil {
		logger.Fatal(err)
	}
}
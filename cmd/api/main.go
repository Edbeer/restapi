package main

import (
	"github.com/Edbeer/restapi/config"
	server "github.com/Edbeer/restapi/internal/transport/rest"
	"github.com/Edbeer/restapi/pkg/db/postgres"
	"github.com/Edbeer/restapi/pkg/db/redis"
	"github.com/Edbeer/restapi/pkg/logger"
)

// @title           restapi
// @version         1.0
// @description     This is an example of an implementation RESTApi

// @BasePath /api/

func main() {
	cfg := config.GetConfig()
	logger := logger.NewApiLogger(cfg)
	logger.InitLogger()

	// postgresql
	psqlClient, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		logger.Fatalf("Postgresql init: %s", err)
	} else {
		logger.Infof("Postgres connected, Status: %#v", psqlClient.Stats())
	}
	defer psqlClient.Close()

	// redis
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	logger.Info("Redis connetcted")

	logger.Info("Starting auth server")
	s := server.NewServer(cfg, psqlClient, redisClient, logger)
	if err := s.Run(); err != nil {
		logger.Fatal(err)
	}
}
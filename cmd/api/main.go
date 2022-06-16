package main

import (
	"log"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/internal/server"
	"github.com/Edbeer/restapi/pkg/db/postgres"
	"github.com/Edbeer/restapi/pkg/db/redis"
	"github.com/Edbeer/restapi/pkg/logger"
)

func main() {
	cfg := config.GetConfig()
	logger := logger.NewApiLogger(cfg)
	logger.InitLogger()

	// postgresql
	psqlClient, err := postgres.NewPsqlClient(cfg)
	if err != nil {
		// logger.Fatalf("Postgresql init: %s", err)
		log.Fatal(err)
	} else {
		// logger.Infof("Postgres connected, Status: %#v", psqlClient.Stat())
		log.Println(psqlClient.Stat())
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
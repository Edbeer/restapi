package main

import (
	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/pkg/logger"
	server "github.com/Edbeer/restapi/internal/transport/rest"
)

func main() {
	cfg := config.GetConfig()
	logger := logger.NewApiLogger(cfg)

	logger.InitLogger()
	logger.Info("Starting auth server")
	s := server.NewServer(cfg, logger)
	logger.Fatal(s.Run())
}
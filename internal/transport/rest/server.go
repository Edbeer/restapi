package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Edbeer/restapi/config"
	"github.com/Edbeer/restapi/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo       *echo.Echo
	config     *config.Config
	psqlClient *pgxpool.Pool
	logger     logger.Logger
}

// New server constructor
func NewServer(config *config.Config, psqlClient *pgxpool.Pool, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), psqlClient: psqlClient, config: config, logger: logger}
}

// Run server depends on config SSL option
func (s *Server) Run() error {
	if s.config.Server.SSL {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * time.Duration(s.config.Server.ReadTimeout)
		s.echo.Server.WriteTimeout = time.Second * time.Duration(s.config.Server.WriteTimeout)

		go func() {
			s.logger.Infof("Server listening on port: %s", s.config.Server.Port)
			if err := s.echo.StartTLS(s.config.Server.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("error starting TLS server %v", err)
			}
		}()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Shutdown(ctx)
	} else {
		e := echo.New()
		if err := s.MapHandlers(e); err != nil {
			return err
		}

		server := &http.Server{
			Addr:           s.config.Server.Port,
			ReadTimeout:    time.Second * time.Duration(s.config.Server.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(s.config.Server.WriteTimeout),
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			s.logger.Infof("Server is listening on port: %s", s.config.Server.Port)
			if err := e.StartServer(server); err != nil {
				s.logger.Fatalf("Error starting Server: %v", err)
			}
		}()

		// Graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}
}

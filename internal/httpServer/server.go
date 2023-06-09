package httpServer

import (
	"UsersService/config"
	"UsersService/pkg/httpErrorHandler"
	"UsersService/pkg/logger"
	"UsersService/pkg/secure"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	apiLogger *logger.ApiLogger
	cfg       *config.Config
	shield    *secure.Shield
	fiber     *fiber.App
}

func NewServer(cfg *config.Config, apiLogger *logger.ApiLogger, handler *httpErrorHandler.HttpErrorHandler, shield *secure.Shield) *Server {
	return &Server{
		fiber: fiber.New(fiber.Config{
			ErrorHandler:          handler.Handler,
			DisableStartupMessage: true,
		}),
		apiLogger: apiLogger,
		shield:    shield,
		cfg:       cfg,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.fiber, s.apiLogger); err != nil {
		s.apiLogger.Fatalf("Cannot map handlers. Error: {%s}", err)
	}

	//s.apiLogger.Infof("Start server on {host:port - %s:%s}", s.cfg.Server.Host, s.cfg.Server.Port)
	s.apiLogger.Infof("Start server on {host:port - %s:%s}", "localhost", "8080")

	//if err := s.fiber.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)); err != nil {
	//	s.apiLogger.Fatalf("Cannot listen. Error: {%s}", err)
	//}

	if err := s.fiber.Listen(fmt.Sprintf("%s:%s", "localhost", "8080")); err != nil {
		s.apiLogger.Fatalf("Cannot listen. Error: {%s}", err)
	}

	return nil
}

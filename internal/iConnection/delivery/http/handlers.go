package http

import (
	"UsersService/config"
	"UsersService/pkg/logger"
)

type ConnectionHandler struct {
	logger *logger.ApiLogger
	cfg    *config.Config
}

func NewUsersHandlers(cfg *config.Config, logger *logger.ApiLogger) *ConnectionHandler {
	return &ConnectionHandler{cfg: cfg, logger: logger}
}

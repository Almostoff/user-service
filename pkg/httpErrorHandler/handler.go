package httpErrorHandler

import (
	"UsersService/config"
	"UsersService/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type HttpErrorHandler struct {
	showUnknownErrors bool
	logger            logger.Logger
}

func NewErrorHandler(c *config.Config, logger logger.Logger) *HttpErrorHandler {
	return &HttpErrorHandler{
		showUnknownErrors: c.Server.ShowUnknownErrorsInResponse,
		logger:            logger,
	}
}

type responseMsg struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (handler *HttpErrorHandler) Handler(c *fiber.Ctx, err error) error {
	var response responseMsg
	var statusCode int

	response.Success = false
	if strings.Contains(err.Error(), "user not found") {
		response.Message = "Invalid public api key"
		statusCode = fiber.StatusUnauthorized
	} else if strings.Contains(err.Error(), "sql") || strings.Contains(err.Error(), "SQL") {
		response.Message = "Internal Server Error"
		statusCode = fiber.StatusInternalServerError
	} else if strings.Contains(err.Error(), "connection refused") {
		response.Message = "Internal Server Error"
		statusCode = fiber.StatusInternalServerError
	} else if statusCode == 0 {
		response.Message = err.Error()
		statusCode = fiber.StatusInternalServerError
	} else {
		response.Message = err.Error()
		statusCode = fiber.StatusInternalServerError
	}
	return c.Status(statusCode).JSON(response)

}

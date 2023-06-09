package utils

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// Read request body and validate
func ReadRequest(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		return err
	}

	return validate.StructCtx(c.Context(), request)
}

func ParseQuery(query string) map[string]string {
	result := make(map[string]string)

	optionsArray := strings.Split(query, "&")
	for _, option := range optionsArray {
		variable, value, _ := strings.Cut(option, "=")

		result[variable] = value
	}
	return result
}

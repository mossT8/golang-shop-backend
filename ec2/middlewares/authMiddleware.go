package middlewares

import (
	"github.com/gofiber/fiber/v2"
	internalMiddleware "tannar.moss/backend/internal/middlewares"
)

func IsAuthenticated(c *fiber.Ctx) error {
	jwt := c.Cookies("jwt")

	err := internalMiddleware.IsAuthenticated(jwt)
	if err != nil {
		return err
	}

	return c.Next()
}

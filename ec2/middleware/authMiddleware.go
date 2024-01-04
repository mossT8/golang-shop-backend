package middleware

import (
	"github.com/gofiber/fiber/v2"
	internalService "tannar.moss/backend/internal/service"
)

func IsAuthenticated(c *fiber.Ctx, service internalService.Auth) error {
	jwt := c.Cookies("jwt")

	err := service.IsAuthenticated(jwt)
	if err != nil {
		return err
	}

	return c.Next()
}

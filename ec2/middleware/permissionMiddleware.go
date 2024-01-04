package middleware

import (
	"github.com/gofiber/fiber/v2"
	internalService "tannar.moss/backend/internal/service"
)

func IsAuthorized(c *fiber.Ctx, page string, service internalService.Auth) error {
	jwt := c.Cookies("jwt")

	err := service.IsAuthorized(jwt, page)
	if err != nil {
		return err
	}

	return nil
}

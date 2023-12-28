package middleware

import (
	"github.com/gofiber/fiber/v2"
	internalMiddleware "tannar.moss/backend/internal/middleware"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	jwt := c.Cookies("jwt")

	err := internalMiddleware.IsAuthorized(jwt, page)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"tannar.moss/backend/ec2/routes"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")
}

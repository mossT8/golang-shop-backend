package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"tannar.moss/backend/ec2/controller"
	"tannar.moss/backend/ec2/middleware"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/service"
)

func Setup(app *fiber.App, controller controller.InternalPluginController) {
	// todo: remove hardcode databse config and use aws secret
	genericUserConfig := mysql.DatabaseConfig{
		Host:              "localhost",
		Port:              5432,
		RequestTimeout:    30,
		ConnectionTimeout: 10,
		Dialect:           "mysql",
		Database:          "go_admin",
		Username:          "root",
		Password:          "root",
	}
	logger := logger.NewSimpleLogger("DEBUG", false)
	dbConn, err := mysql.NewDbConnection(genericUserConfig, genericUserConfig)
	if err != nil {
		logger.Errorf("Unabled to connect to database: %s", err.Error())
		return
	}

	// setup middleware
	userRepo := repository.NewMySqlUserRepository(logger, *dbConn)
	validatorService := service.NewValidator(logger, *validator.New())
	authService := service.NewAuthService(validatorService, userRepo, logger)

	// auth routes
	app.Post("/api/register", controller.Register())
	app.Put("/api/login", controller.Login())

	app.Use(func(c *fiber.Ctx) error {
		return middleware.IsAuthenticated(c, authService)
	})

	// management routes
	app.Put("/api/users/info", controller.UpdateInfo())
	app.Put("/api/users/password", controller.UpdatePassword())

	app.Get("/api/user", controller.User())
	app.Post("/api/logout", controller.Logout())

	app.Get("/api/users", controller.AllUsers())
	app.Get("/api/users/:id", controller.GetUser())
	app.Post("/api/users", controller.CreateUser())
	app.Put("/api/users/:id", controller.UpdateUser())
	app.Delete("/api/users/:id", controller.DeleteUser())

	app.Get("/api/roles", controller.AllRoles())
	app.Get("/api/roles/:id", controller.GetRole())
	app.Post("/api/roles", controller.CreateRole())
	app.Put("/api/roles/:id", controller.UpdateRole())
	app.Delete("/api/roles/:id", controller.DeleteRole())

	app.Get("/api/permissions", controller.AllPermissions())

	app.Get("/api/products", controller.AllProducts())
	app.Get("/api/products/:id", controller.GetProduct())
	app.Post("/api/products", controller.CreateProduct())
	app.Put("/api/products/:id", controller.UpdateProduct())
	app.Delete("/api/products/:id", controller.DeleteProduct())

	app.Post("/api/upload", controller.Upload())
	app.Static("/api/uploads", "/uploads")

	app.Get("/api/chart", controller.Chart())

	// orders route
	app.Get("/api/orders", controller.AllOrders())
	app.Get("/api/order/:id", controller.GetOrder())
	app.Post("/api/export", controller.Export())
	app.Put("/api/order/:id", controller.UpdateOrder())
	app.Post("/api/order", controller.CreateOrder())
}

package controller

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"tannar.moss/backend/internal/constant"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/service"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

type InternalPluginController interface {
	GetPublicService() service.Public
	Register() fiber.Handler
	Login() fiber.Handler
	User() fiber.Handler
	Logout() fiber.Handler
	UpdatePassword() fiber.Handler
	UpdateInfo() fiber.Handler
	Upload() fiber.Handler
	AllOrders() fiber.Handler
	AddOrder() fiber.Handler
	CreateOrder() fiber.Handler
	UpdateOrder() fiber.Handler
	GetOrder() fiber.Handler
	DeleteOrder() fiber.Handler
	Export() fiber.Handler
	CreateFile() fiber.Handler
	AllPermissions() fiber.Handler
	AllProducts() fiber.Handler
	CreateProduct() fiber.Handler
	GetProduct() fiber.Handler
	UpdateProduct() fiber.Handler
	DeleteProduct() fiber.Handler
	AllRoles() fiber.Handler
	CreateRole() fiber.Handler
	UpdateRole() fiber.Handler
	DeleteRole() fiber.Handler
	GetRole() fiber.Handler
	AllUsers() fiber.Handler
	CreateUser() fiber.Handler
	GetUser() fiber.Handler
	UpdateUser() fiber.Handler
	DeleteUser() fiber.Handler
	Chart() fiber.Handler
}

type InternalPluginControllerImpl struct {
	publicService  service.Public
	privateService service.Private
	logger         logger.Logger
}

func (controller *InternalPluginControllerImpl) GetPublicService() service.Public {
	return controller.publicService
}

func (controller *InternalPluginControllerImpl) getUserIdSession(context *fiber.Ctx) (uint64, error) {
	jwt, err := controller.getJwtTokenFromSession(context)
	if err != nil {
		return 0, err
	}
	issuer, err := utils.GetIssuerFromJwt(jwt, constant.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		controller.logger.Errorf("Cant get issuer from jwt '%s'", issuer)
		return 0, types.NewInternalServerError()
	}
	userId, err := strconv.ParseUint(issuer, 10, 64)
	if err != nil {
		controller.logger.Errorf("Cant parse as Uint: '%s'", issuer)
		return 0, types.NewInternalServerError()
	}
	return userId, nil
}

func (controller *InternalPluginControllerImpl) getJwtTokenFromSession(context *fiber.Ctx) (string, error) {
	authHeader := context.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	cookie := context.Cookies("jwt")
	if cookie != "" {
		return cookie, nil
	}

	return "", errors.New("jwt token not found in headers or cookies")
}

// AddOrder implements InternalPluginController.
func (InternalPluginControllerImpl) AddOrder() fiber.Handler {
	panic("unimplemented")
}

// GetOrder implements InternalPluginController.
func (InternalPluginControllerImpl) GetOrder() fiber.Handler {
	panic("unimplemented")
}

// DeleteOrder implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteOrder() fiber.Handler {
	panic("unimplemented")
}

// UpdateOrder implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateOrder() fiber.Handler {
	panic("unimplemented")
}

// CreateOrder implements InternalPluginController.
func (InternalPluginControllerImpl) CreateOrder() fiber.Handler {
	panic("unimplemented")
}

// AllOrders implements InternalPluginController.
func (InternalPluginControllerImpl) AllOrders() fiber.Handler {
	panic("unimplemented")
}

// AllPermissions implements InternalPluginController.
func (InternalPluginControllerImpl) AllPermissions() fiber.Handler {
	panic("unimplemented")
}

// AllProducts implements InternalPluginController.
func (InternalPluginControllerImpl) AllProducts() fiber.Handler {
	panic("unimplemented")
}

// AllRoles implements InternalPluginController.
func (InternalPluginControllerImpl) AllRoles() fiber.Handler {
	panic("unimplemented")
}

// AllUsers implements InternalPluginController.
func (InternalPluginControllerImpl) AllUsers() fiber.Handler {
	panic("unimplemented")
}

// Chart implements InternalPluginController.
func (InternalPluginControllerImpl) Chart() fiber.Handler {
	panic("unimplemented")
}

// CreateFile implements InternalPluginController.
func (InternalPluginControllerImpl) CreateFile() fiber.Handler {
	panic("unimplemented")
}

// CreateProduct implements InternalPluginController.
func (InternalPluginControllerImpl) CreateProduct() fiber.Handler {
	panic("unimplemented")
}

// CreateRole implements InternalPluginController.
func (InternalPluginControllerImpl) CreateRole() fiber.Handler {
	panic("unimplemented")
}

// CreateUser implements InternalPluginController.
func (InternalPluginControllerImpl) CreateUser() fiber.Handler {
	panic("unimplemented")
}

// DeleteProduct implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteProduct() fiber.Handler {
	panic("unimplemented")
}

// DeleteRole implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteRole() fiber.Handler {
	panic("unimplemented")
}

// DeleteUser implements InternalPluginController.
func (InternalPluginControllerImpl) DeleteUser() fiber.Handler {
	panic("unimplemented")
}

// Export implements InternalPluginController.
func (InternalPluginControllerImpl) Export() fiber.Handler {
	panic("unimplemented")
}

// GetProduct implements InternalPluginController.
func (InternalPluginControllerImpl) GetProduct() fiber.Handler {
	panic("unimplemented")
}

// GetRole implements InternalPluginController.
func (InternalPluginControllerImpl) GetRole() fiber.Handler {
	panic("unimplemented")
}

// GetUser implements InternalPluginController.
func (InternalPluginControllerImpl) GetUser() fiber.Handler {
	panic("unimplemented")
}

func (controller *InternalPluginControllerImpl) marshalErrorResponse(context *fiber.Ctx, err error) error {
	typedErr, ok := err.(*types.SocketError)
	if !ok {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return context.Status(typedErr.StatusCode()).JSON(fiber.Map{
		"message": typedErr.Error(),
	})
}

// Login implements InternalPluginController.
func (controller *InternalPluginControllerImpl) Login() fiber.Handler {
	return func(context *fiber.Ctx) error {
		loginResponse, err := controller.publicService.Login(string(context.Body()))

		if err != nil {
			return controller.marshalErrorResponse(context, err)
		}
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    loginResponse.Jwt,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		}
		context.Cookie(&cookie)
		return context.JSON(loginResponse)
	}
}

// Logout implements InternalPluginController.
func (InternalPluginControllerImpl) Logout() fiber.Handler {
	panic("unimplemented")
}

func (controller InternalPluginControllerImpl) Register() fiber.Handler {
	return func(context *fiber.Ctx) error {
		loginResponse, err := controller.publicService.Register(string(context.Body()))

		if err != nil {
			return controller.marshalErrorResponse(context, err)
		}
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    loginResponse.Jwt,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		}
		context.Cookie(&cookie)
		return context.JSON(loginResponse)
	}
}

func (controller *InternalPluginControllerImpl) UpdateInfo() fiber.Handler {
	return func(context *fiber.Ctx) error {
		userId, err := controller.getUserIdSession(context)
		if err != nil {
			return controller.marshalErrorResponse(context, err)
		}
		userResponse, err := controller.privateService.UpdateUserInfo(userId, string(context.Body()), userId)
		if err != nil {
			return context.JSON(utils.FormatErrorAPIGatewayResponse(err))
		}
		return context.JSON(userResponse)
	}
}

func (controller *InternalPluginControllerImpl) UpdatePassword() fiber.Handler {
	return func(context *fiber.Ctx) error {
		userId, err := controller.getUserIdSession(context)
		if err != nil {
			return controller.marshalErrorResponse(context, err)
		}
		updatedResponse, err := controller.privateService.UpdateUserPassword(userId, string(context.Body()), userId)
		if err != nil {
			return controller.marshalErrorResponse(context, err)
		}
		return context.JSON(updatedResponse)
	}
}

// UpdateProduct implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateProduct() fiber.Handler {
	panic("unimplemented")
}

// UpdateRole implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateRole() fiber.Handler {
	panic("unimplemented")
}

// UpdateUser implements InternalPluginController.
func (InternalPluginControllerImpl) UpdateUser() fiber.Handler {
	panic("unimplemented")
}

// Upload implements InternalPluginController.
func (InternalPluginControllerImpl) Upload() fiber.Handler {
	panic("unimplemented")
}

// User implements InternalPluginController.
func (InternalPluginControllerImpl) User() fiber.Handler {
	panic("unimplemented")
}

func NewInternalPluginController() InternalPluginController {
	// todo: remove hardcode databse config and use aws secret
	genericUserConfig := mysql.DatabaseConfig{
		Host:              "localhost",
		Port:              3306,
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
		panic("DB down!!!")
	}

	userRepo := repository.NewMySqlUserRepository(logger, *dbConn)
	validatorService := service.NewValidator(logger, *validator.New())
	publicService := service.NewPublicService(validatorService, userRepo, logger)
	privateService := service.NewPrivateService(validatorService, userRepo, logger)

	logger.Info("System started... ")

	return &InternalPluginControllerImpl{
		publicService:  publicService,
		privateService: privateService,
		logger:         logger,
	}
}

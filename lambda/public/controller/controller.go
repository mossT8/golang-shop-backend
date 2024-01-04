package controller

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/service"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/lambda"
	"tannar.moss/backend/lambda/public/model"
)

type PublicController struct {
	Service service.Auth
	Logger  logger.Logger
}

func NewPrivateController() (lambda.Controller, error) {
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

	dbConn, err := mysql.NewDbConnection(genericUserConfig, genericUserConfig)
	if err != nil {
		return nil, types.NewInternalServerError()
	}

	logger := logger.NewSimpleLogger("DEBUG", false)
	userRepo := repository.NewMySqlUserRepository(logger, *dbConn)
	validatorService := service.NewValidator(logger, *validator.New())
	authService := service.NewAuthService(validatorService, userRepo, logger)

	return &PublicController{
		Service: authService,
		Logger:  logger,
	}, nil
}

// PostProcess implements lambda.Controller.
func (*PublicController) PostProcess(request model.Response) (string, error) {
	panic("unimplemented")
}

// PreProcess implements lambda.Controller.
func (*PublicController) PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (string, string, string, error) {
	panic("unimplemented")
}

// Process implements lambda.Controller.
func (*PublicController) Process(requestType string, path string, body string) (*model.Response, error) {
	panic("unimplemented")
}

// PublishLogs implements lambda.Controller.
func (*PublicController) PublishLogs() {
	panic("unimplemented")
}

// Shutdown implements lambda.Controller.
func (*PublicController) Shutdown() {
	panic("unimplemented")
}

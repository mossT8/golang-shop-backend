package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

const (
	PUT  = "PUT"
	POST = "POST"
)

func NewPublicController(logLevel string, publishLogs bool) (lambda.Controller, error) {
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

	dbConn, err := mysql.NewDbConnection(genericUserConfig, genericUserConfig)
	if err != nil {
		return nil, types.NewInternalServerError()
	}

	logger := logger.NewSimpleLogger(logLevel, publishLogs)
	userRepo := repository.NewMySqlUserRepository(logger, *dbConn)
	validatorService := service.NewValidator(logger, *validator.New())
	authService := service.NewAuthService(validatorService, userRepo, logger)

	return &PublicController{
		Service: authService,
		Logger:  logger,
	}, nil
}

func (this *PublicController) PostProcess(response model.Response) (string, error) {
	responseString, err := json.Marshal(response)
	if err != nil {
		this.Logger.Errorf("Error converting response to string: %v", err)
		return "", err
	}

	this.Logger.Infof("Response: %s", responseString)
	return string(responseString), nil
}

func (this *PublicController) PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (string, string, string, error) {
	this.Logger.SetTraceId(uuid.NewString())
	return event.HTTPMethod, event.Path, event.Body, nil
}

func (this *PublicController) Process(requestType string, path string, body string) (*model.Response, error) {
	switch requestType {
	case "POST":
		return this.handlePostRequest(path, body)
	case "PUT":
		return this.handlePutRequest(path, body)
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (this *PublicController) handlePostRequest(path string, body string) (*model.Response, error) {
	switch path {
	case "/api/register":
		loginResponse, err := this.Service.Register(body)
		if err != nil {
			return nil, err
		}
		return &model.Response{
			LoginResponse: *loginResponse,
		}, nil
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (this *PublicController) handlePutRequest(path string, body string) (*model.Response, error) {
	switch path {
	case "/api/login":
		loginResponse, err := this.Service.Login(body)
		if err != nil {
			return nil, err
		}

		return &model.Response{
			LoginResponse: *loginResponse,
		}, nil
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (this *PublicController) PublishLogs() {
	this.Logger.Info("To do add SQS for sumo logging")
}

func (this *PublicController) Shutdown() {
	this.Service.Shutdown()
	this = nil
}

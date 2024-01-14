package controller

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"tannar.moss/backend/internal/constant"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/service"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/lambda/public/model"
)

type Controller interface {
	PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (string, string, string, error)
	Process(requestType string, path string, body string) (*model.Response, error)
	PostProcess(response model.Response) (string, error)
	PublishLogs()
	Shutdown()
}

type PublicController struct {
	Service service.Public
	Logger  logger.Logger
}

func NewPublicController(logLevel string, publishLogs bool) (Controller, error) {
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
	publicService := service.NewPublicService(validatorService, userRepo, logger)

	return &PublicController{
		Service: publicService,
		Logger:  logger,
	}, nil
}

func (c *PublicController) PostProcess(response model.Response) (string, error) {
	responseString, err := json.Marshal(response)
	if err != nil {
		c.Logger.Errorf("Error converting response to string: %v", err)
		return "", err
	}

	c.Logger.Infof("Response: %s", responseString)
	return string(responseString), nil
}

func (c *PublicController) PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (string, string, string, error) {
	c.Logger.SetTraceId(uuid.NewString())
	return event.HTTPMethod, event.Path, event.Body, nil
}

func (c *PublicController) Process(requestType string, path string, body string) (*model.Response, error) {
	switch requestType {
	case constant.POST:
		return c.handlePostRequest(path, body)
	case constant.PUT:
		return c.handlePutRequest(path, body)
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (c *PublicController) handlePostRequest(path string, body string) (*model.Response, error) {
	switch path {
	case "/api/register":
		loginResponse, err := c.Service.Register(body)
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

func (c *PublicController) handlePutRequest(path string, body string) (*model.Response, error) {
	switch path {
	case "/api/login":
		loginResponse, err := c.Service.Login(body)
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

func (c *PublicController) PublishLogs() {
	c.Logger.PublishSumoLogs()
}

func (c *PublicController) Shutdown() {
	c.Service.Shutdown()
	c = nil
}

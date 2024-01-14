package controller

import (
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"tannar.moss/backend/internal/constant"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/service"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
	"tannar.moss/backend/lambda/private/model"
)

type Controller interface {
	PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (uint64, string, string, string, error)
	Process(signInUserId uint64, requestType string, path string, body string) (*model.Response, error)
	PostProcess(response model.Response) (string, error)
	PublishLogs()
	Shutdown()
}

type PrivateController struct {
	service service.Private
	logger  logger.Logger
}

func (c *PrivateController) retrieveUserIdFromJWT(jwt string) (uint64, error) {
	issuer, err := utils.GetIssuerFromJwt(jwt, constant.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		c.logger.Errorf("Cant get issuer from jwt '%s'", issuer)
		return 0, types.NewInternalServerError()
	}
	userId, err := strconv.ParseUint(issuer, 10, 64)
	if err != nil {
		c.logger.Errorf("Cant parse as Uint: '%s'", issuer)
		return 0, types.NewInternalServerError()
	}

	return userId, nil
}
func NewPrivateController(logLevel string, publishLogs bool) (Controller, error) {
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
	PrivateService := service.NewPrivateService(validatorService, userRepo, logger)

	return &PrivateController{
		service: PrivateService,
		logger:  logger,
	}, nil
}

func (c *PrivateController) PostProcess(response model.Response) (string, error) {
	responseString, err := json.Marshal(response)
	if err != nil {
		c.logger.Errorf("Error converting response to string: %v", err)
		return "", err
	}

	c.logger.Infof("Response: %s", responseString)
	return string(responseString), nil
}

func (c *PrivateController) PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (uint64, string, string, string, error) {
	c.logger.SetTraceId(uuid.NewString())
	jwtToken := event.Headers["Authorization"]
	userId, err := c.retrieveUserIdFromJWT(jwtToken)
	return userId, event.HTTPMethod, event.Path, event.Body, err
}

func (c *PrivateController) Process(userId uint64, requestType string, path string, body string) (*model.Response, error) {
	switch requestType {
	case constant.GET:
		return c.handleGetRequest(userId, path, body)
	case constant.POST:
		return c.handlePostRequest(userId, path, body)
	case constant.PUT:
		return c.handlePutRequest(userId, path, body)
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (c *PrivateController) handleGetRequest(userId uint64, path string, body string) (*model.Response, error) {
	switch path {
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (c *PrivateController) handlePostRequest(userId uint64, path string, body string) (*model.Response, error) {
	switch path {
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (c *PrivateController) handlePutRequest(userId uint64, path string, body string) (*model.Response, error) {
	switch path {
	case "/api/users/info":
		userResponse, err := c.service.UpdateUserInfo(userId, body, userId)
		if err != nil {
			return nil, err
		}

		return &model.Response{
			User: *userResponse,
		}, nil
	case "/api/users/password":
		userResponse, err := c.service.UpdateUserPassword(userId, body, userId)
		if err != nil {
			return nil, err
		}

		return &model.Response{
			User: *userResponse,
		}, nil
	default:
		return nil, types.NewNotImplementedError()
	}
}

func (c *PrivateController) PublishLogs() {
	c.logger.PublishSumoLogs()
}

func (c *PrivateController) Shutdown() {
	c.service.Shutdown()
	c = nil
}

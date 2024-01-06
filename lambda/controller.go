package lambda

import (
	"github.com/aws/aws-lambda-go/events"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository/mysql"
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

type PrivateController struct {
	Logger logger.Logger
}

func NewPrivateController() (Controller, error) {
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

	_, err := mysql.NewDbConnection(genericUserConfig, genericUserConfig)
	if err != nil {
		return nil, types.NewInternalServerError()
	}

	logger := logger.NewSimpleLogger("DEBUG", false)

	return &PrivateController{
		Logger: logger,
	}, nil
}

// PostProcess implements Controller.
func (*PrivateController) PostProcess(request model.Response) (string, error) {
	panic("unimplemented")
}

// PreProcess implements Controller.
func (*PrivateController) PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (string, string, string, error) {
	panic("unimplemented")
}

// Process implements Controller.
func (*PrivateController) Process(requestType string, path string, body string) (*model.Response, error) {
	panic("unimplemented")
}

// PublishLogs implements Controller.
func (*PrivateController) PublishLogs() {
	panic("unimplemented")
}

// Shutdown implements Controller.
func (*PrivateController) Shutdown() {
	panic("unimplemented")
}

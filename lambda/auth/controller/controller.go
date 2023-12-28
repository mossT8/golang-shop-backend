package controller

import (
	"github.com/aws/aws-lambda-go/events"
	"tannar.moss/backend/lambda/auth/model"
)

type Controller interface {
	PreProcess(event events.APIGatewayWebsocketProxyRequest, loglevel string, pushLogs bool) (string, string, string, error)
	Process(requestType string, path string, body string) (*model.Response, error)
	PostProcess(request model.Response) (string, error)
	PublishLogs()
	Shutdown()
}

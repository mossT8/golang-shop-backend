package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"tannar.moss/backend/internal/utils"
)

var invokeCount = 0
var MAX_INVOKE = 0

func handlerEvent(_ context.Context, event events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	invokeCount++
	if invokeCount >= MAX_INVOKE {
		// restart
	}
	response := processEvent(event)

	return response, nil
}

func processEvent(event events.APIGatewayWebsocketProxyRequest) *events.APIGatewayProxyResponse {

	return utils.FormatGatewayResponse(200, "Hello World!")
}

func main() {
	lambda.Start(handlerEvent)
}

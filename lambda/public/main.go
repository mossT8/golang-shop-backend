package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"tannar.moss/backend/internal/utils"
	internalLambda "tannar.moss/backend/lambda"
	"tannar.moss/backend/lambda/public/controller"
)

var invokeCount = 0
var lambdaController internalLambda.Controller

func handlerEvent(_ context.Context, event events.APIGatewayWebsocketProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logLevel := utils.Getenv("LOG_LEVEL", "INFO")
	pushLogs := utils.SafeBool(os.Getenv("PUSH_LOGS"), false)

	if invokeCount >= utils.SafeAtoi(os.Getenv("MAX_INVOKE"), 15) {
		lambdaController.Shutdown()
		invokeCount = 0
	}

	var err error
	if invokeCount == 0 {
		lambdaController, err = controller.NewPublicController(logLevel, pushLogs)

		if err != nil {
			return utils.FormatErrorAPIGatewayResponse(err), nil
		}
	}

	invokeCount++
	response := processEvent(event, logLevel, pushLogs)
	lambdaController.PublishLogs()

	return response, nil
}

func processEvent(event events.APIGatewayWebsocketProxyRequest, logLevel string, publishLogs bool) *events.APIGatewayProxyResponse {
	httpType, path, body, err := lambdaController.PreProcess(event, logLevel, publishLogs)
	if err != nil {
		return utils.FormatErrorAPIGatewayResponse(err)
	}

	response, err := lambdaController.Process(httpType, path, body)
	if err != nil {
		return utils.FormatErrorAPIGatewayResponse(err)
	}

	processedResponse, err := lambdaController.PostProcess(*response)
	if err != nil {
		return utils.FormatErrorAPIGatewayResponse(err)
	}

	return utils.FormatGatewayResponse(http.StatusOK, processedResponse)
}

func main() {
	lambda.Start(handlerEvent)
}

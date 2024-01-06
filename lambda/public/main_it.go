package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestLambdaHandler_withPostRequest_shouldSuccessfullProcess(t *testing.T) {

	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("PUSH_TO_SUMO", "false")

	event := events.APIGatewayWebsocketProxyRequest{
		HTTPMethod: "POST",
		Path:       "/api/register",
		Body: `{
			"first_name": "John",
			"last_name": "Doe",
			"email": "johndoe@example.com",
			"password": "password123",
			"confirm_password": "password123"
		}`,
	}

	response, err := handlerEvent(nil, event)
	if err != nil {
		t.Errorf("LambdaHandler() error = %v", err)
	} else {
		t.Log(response.Body)
	}
}

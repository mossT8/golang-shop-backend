package service

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/types"
)

type Validator interface {
	MarshalAndValidateLoginRequest(body string, request any) error
}

type SimpleValidator struct {
	Logger   logger.Logger
	Validate validator.Validate
}

func NewValidator(logger logger.Logger, validate validator.Validate) Validator {
	return &SimpleValidator{
		Logger:   logger,
		Validate: validate,
	}
}

func getStuctName(request any) string {
	structType := reflect.TypeOf(request)
	if structType.Kind() == reflect.Struct {
		return structType.Name()
	} else {
		return ""
	}
}

func (this *SimpleValidator) MarshalAndValidateLoginRequest(body string, request any) error {
	err := json.Unmarshal([]byte(body), &request)
	if err != nil {
		this.Logger.Infof("Error marshaling %s with body '%s' due to '%s'", getStuctName(request), body, err.Error())
		return types.NewInvalidInputError()
	}

	err = this.Validate.Struct(request)
	if err != nil {
		this.Logger.Infof("Error validating %s with body '%s' due to '%s'", getStuctName(request), body, err.Error())
		return types.NewInvalidInputError()
	}

	return nil
}

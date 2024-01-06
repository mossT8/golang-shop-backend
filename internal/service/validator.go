package service

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/model"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

type Validator interface {
	MarshalAndValidateLoginRequest(body string) (*model.LoginRequest, error)
	MarshalAndValidateRegisterRequest(body string) (*model.UserRequest, error)
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

func (this *SimpleValidator) MarshalAndValidateRegisterRequest(body string) (*model.UserRequest, error) {
	body = utils.FormatJSONString(body)
	var request *model.UserRequest
	err := json.Unmarshal([]byte(body), &request)
	if err != nil {
		this.Logger.Infof("Error marshaling request with body '%s' due to '%s'", body, err.Error())
		return nil, types.NewInvalidInputError()
	}

	err = this.Validate.Struct(request)
	if err != nil {
		this.Logger.Infof("Error validating %s with body '%s' due to '%s'", getStuctName(request), body, err.Error())
		return nil, types.NewInvalidInputError()
	}

	return request, nil
}

func (this *SimpleValidator) MarshalAndValidateLoginRequest(body string) (*model.LoginRequest, error) {
	body = utils.FormatJSONString(body)
	var request *model.LoginRequest
	err := json.Unmarshal([]byte(body), &request)
	if err != nil {
		this.Logger.Infof("Error marshaling request with body '%s' due to '%s'", body, err.Error())
		return nil, types.NewInvalidInputError()
	}

	err = this.Validate.Struct(request)
	if err != nil {
		this.Logger.Infof("Error validating %s with body '%s' due to '%s'", getStuctName(request), body, err.Error())
		return nil, types.NewInvalidInputError()
	}

	return request, nil
}

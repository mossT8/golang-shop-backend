package service

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

type Validator interface {
	MarshalAndValidateREQ(body string, request any) error
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
func (validator *SimpleValidator) MarshalAndValidateREQ(body string, request any) error {
	body = utils.FormatJSONString(body)
	err := json.Unmarshal([]byte(body), &request)
	if err != nil {
		validator.Logger.Infof("Error marshaling request: %s", err.Error())
		return types.NewInvalidInputError()
	}

	err = validator.Validate.Struct(request)
	if err != nil {
		validator.Logger.Infof("Error validating: %s", err.Error())
		return types.NewInvalidInputError()
	}

	return nil
}

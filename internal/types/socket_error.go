package types

import "tannar.moss/backend/internal/constant"

type SocketError struct {
	statusCode int
	message    string
}

func NewSocketError(statusCode int, message string) *SocketError {
	return &SocketError{
		statusCode: statusCode,
		message:    message,
	}
}

func (e *SocketError) Error() string {
	return e.message
}

func (e *SocketError) StatusCode() int {
	return e.statusCode
}

func NewInternalServerError() error {
	return NewSocketError(constant.InternalServerErrorCode, constant.InternalServerErrorName)
}

func NewBadRequestError() error {
	return NewSocketError(constant.BadRequestCode, constant.BadRequestName)
}

func NewInvalidInputError() error {
	return NewSocketError(constant.InvalidInputCode, constant.InvalidInputErrorName)
}

func NewNotImplementedError() error {
	return NewSocketError(constant.NotImplementedCode, constant.NotImplementedErrorName)
}

func NewNoTFoundOrNoRecordError() error {
	return NewSocketError(constant.NotFoundCode, constant.NotFoundErrorName)
}

func NewUnauthorizedError() error {
	return NewSocketError(constant.UnauthorizedCode, constant.UnauthorizedRequestName)
}

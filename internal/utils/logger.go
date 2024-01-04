package utils

import (
	"fmt"

	"tannar.moss/backend/internal/logger"
)

func LogPreparingError(queryName string, logger logger.Logger, err error) {
	logger.Error(fmt.Sprintf("Error preparing '%s' query: %s", queryName, err.Error()))
}

func LogExecutingError(queryName string, logger logger.Logger, err error) {
	logger.Error(fmt.Sprintf("Error executing '%s' query: %s", queryName, err.Error()))
}

func LogBeginingTnxError(queryName string, logger logger.Logger, err error) {
	logger.Error(fmt.Sprintf("Error begining transaction '%s' query: %s", queryName, err.Error()))
}

func LogCommitError(queryName string, logger logger.Logger, err error) {
	logger.Error(fmt.Sprintf("Error commiting transaction '%s' query: %s", queryName, err.Error()))
}

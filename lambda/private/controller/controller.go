package controller

import (
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/service"
)

type PrivateController struct {
	Service service.Auth
	Logger  logger.Logger
}

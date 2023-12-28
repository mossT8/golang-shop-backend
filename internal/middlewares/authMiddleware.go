package middlewares

import (
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

func IsAuthenticated(jwt string) error {
	if _, err := utils.ParseJwt(jwt, "secret"); err != nil {
		return types.NewNotImplementedError()
	}

	return nil
}

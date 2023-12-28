package middleware

import (
	"fmt"

	"tannar.moss/backend/internal/utils"
)

func IsAuthorized(jwt string, page string) error {
	Id, err := utils.ParseJwt(jwt, "secret")

	fmt.Println(Id)

	if err != nil {
		return err
	}

	// check page permission

	return nil
}

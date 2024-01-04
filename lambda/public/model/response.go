package model

import internalModel "tannar.moss/backend/internal/model"

type Response struct {
	LoginResponse internalModel.LoginResponse
	User          internalModel.UserResponse
}

package service

import (
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/model"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/types"
)

type Private interface {
	UpdateUserInfo(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error)
	UpdateUserPassword(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error)
	Shutdown()
}

type PrivateService struct {
	validator Validator
	userRepo  repository.UserRepository
	logger    logger.Logger
}

func (p *PrivateService) UpdateUserInfo(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error) {
	var updateUserRequest model.UserUpdateRequest
	err := p.validator.MarshalAndValidateREQ(body, &updateUserRequest)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.Update(userId, updateUserRequest.FirstName, updateUserRequest.LastName, updatingUserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PrivateService) UpdateUserPassword(userId uint64, body string, updatingUserId uint64) (*model.UserResponse, error) {
	var updateUserRequest model.ChangePasswordRequest
	err := p.validator.MarshalAndValidateREQ(body, &updateUserRequest)
	if err != nil {
		return nil, err
	}

	if updateUserRequest.ConfirmPassowrd != updateUserRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	user, err := p.userRepo.ResetPassword(userId, updateUserRequest.Password, updatingUserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PrivateService) Shutdown() {
	p.userRepo.Shutdown()
}

func NewPrivateService(validator Validator, userRepo repository.UserRepository, logger logger.Logger) Private {
	return &PrivateService{
		validator: validator,
		userRepo:  userRepo,
		logger:    logger,
	}
}

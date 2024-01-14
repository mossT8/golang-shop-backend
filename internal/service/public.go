package service

import (
	"strconv"

	"tannar.moss/backend/internal/constant"
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/model"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

type Public interface {
	IsAuthenticated(jwt string) error
	IsAuthorized(jwt string, page string) error
	User(jwt string) (*model.UserResponse, error)
	Logout(jwt string) error
	Register(body string) (*model.LoginResponse, error)
	Login(body string) (*model.LoginResponse, error)
	Shutdown()
}

type PublicService struct {
	validator Validator
	userRepo  repository.UserRepository
	logger    logger.Logger
}

func (auth *PublicService) Shutdown() {
	auth.userRepo.Shutdown()
	auth = nil
}

func NewPublicService(validator Validator, userReo repository.UserRepository, logger logger.Logger) Public {
	return &PublicService{
		validator: validator,
		userRepo:  userReo,
		logger:    logger,
	}
}

func (auth *PublicService) generateLoginResponseFromUser(user model.UserResponse) (*model.LoginResponse, error) {

	token, expireAt, err := utils.GenerateJwt(utils.UintToString(user.ID), constant.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		auth.logger.Infof("Cant generate Jwt for userId = '%d' due to '%s'", user.ID, err.Error())
		return nil, types.NewInternalServerError()
	}
	auth.logger.Debugf("User logged in '%d' with token '%s'", user.ID, token)

	return &model.LoginResponse{
		Jwt:      token,
		ExpireAt: expireAt,
	}, nil
}

func (auth *PublicService) checkJwt(jwt string) (string, error) {
	Id, err := utils.ParseJwt(jwt, constant.PASSWORD_SECRET_HASHING_KEY)

	if err != nil {
		auth.logger.Infof("Token Unabled to Parse: '%s'", err.Error())
		return "", err
	}

	auth.logger.Debugf("Token Parsed with ID: '%s'", Id)

	return Id, nil
}

func (auth *PublicService) IsAuthenticated(jwt string) error {
	if _, err := auth.checkJwt(jwt); err != nil {
		return err
	}

	return nil
}

func (auth *PublicService) IsAuthorized(jwt string, page string) error {
	_, err := auth.checkJwt(jwt)
	if err != nil {
		return err
	}

	// check page permission

	return nil
}

func (auth PublicService) Login(body string) (*model.LoginResponse, error) {
	var loginRequest model.LoginRequest
	err := auth.validator.MarshalAndValidateREQ(body, &loginRequest)
	if err != nil {
		return nil, err
	}

	user, err := auth.userRepo.GetByEmail(loginRequest.Username)
	if err != nil {
		return nil, err
	}

	if !utils.ComparePassword(user.HashedPassword, loginRequest.Password) {
		auth.logger.Infof("Failed login attempt for '%s' with '%s'", loginRequest.Username, loginRequest.Password)
		return nil, types.NewUnauthorizedError()
	}

	return auth.generateLoginResponseFromUser(*user)

}

func (auth *PublicService) Logout(jwt string) error {
	userId, err := utils.GetIssuerFromJwt(jwt, constant.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		auth.logger.Errorf("Unabled to logout token: '%s'", jwt)
	}
	auth.logger.Debugf("Logged out user with Id = %d and token: '%s'", userId, jwt)
	return nil
}

func (auth *PublicService) Register(body string) (*model.LoginResponse, error) {
	var registerRequest model.UserRequest
	err := auth.validator.MarshalAndValidateREQ(body, &registerRequest)
	if err != nil {
		return nil, err
	}

	if registerRequest.ConfirmPassword != registerRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	user, err := auth.userRepo.Register(registerRequest.FirstName, registerRequest.LastName, registerRequest.Email, registerRequest.Password, constant.CUSTOMER_ROLE_ID)
	if err != nil {
		return nil, err
	}

	return auth.generateLoginResponseFromUser(*user)
}

func (auth *PublicService) User(jwt string) (*model.UserResponse, error) {
	issuer, err := utils.GetIssuerFromJwt(jwt, constant.PASSWORD_SECRET_HASHING_KEY)
	if err != nil {
		auth.logger.Errorf("Cant get issuer from jwt '%s'", issuer)
		return nil, types.NewInternalServerError()
	}
	userId, err := strconv.ParseUint(issuer, 10, 64)
	if err != nil {
		auth.logger.Errorf("Cant parse as Uint: '%s'", issuer)
		return nil, types.NewInternalServerError()
	}
	user, err := auth.userRepo.GetByID(userId)
	if err != nil {
		auth.logger.Errorf("No user from token: '%s'", jwt)
		return nil, types.NewInternalServerError()
	}
	return user, nil
}

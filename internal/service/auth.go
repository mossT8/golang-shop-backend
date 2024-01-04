package service

import (
	"strconv"

	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/model"
	"tannar.moss/backend/internal/repository"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

type Auth interface {
	IsAuthenticated(jwt string) error
	IsAuthorized(jwt string, page string) error
	User(jwt string) (*model.UserResponse, error)
	Logout(jwt string) error
	Register(body string) (*model.LoginResponse, error)
	Login(body string) (*model.LoginResponse, error)
}

type SimpleAuth struct {
	Validator Validator
	UserRepo  repository.UserRepository
	Logger    logger.Logger
}

const SECRET = "SECRET"

func NewAuthService(validator Validator, userReo repository.UserRepository, logger logger.Logger) Auth {
	return &SimpleAuth{
		Validator: validator,
		UserRepo:  userReo,
		Logger:    logger,
	}
}

func (this *SimpleAuth) generateLoginResponseFromUser(user model.UserResponse) (*model.LoginResponse, error) {

	token, expireAt, err := utils.GenerateJwt(utils.UintToString(user.ID), SECRET)
	if err != nil {
		this.Logger.Infof("UNABLE TO GENERATE JWT FOR USER '%d' REASON '%s'", user.ID, err.Error())
		return nil, types.NewInternalServerError()
	}
	this.Logger.Debugf("SUCCESUFUL LOGIN FOR USER '%d' TOKEN '%s'", user.ID, token)

	return &model.LoginResponse{
		Jwt:      token,
		ExpireAt: expireAt,
	}, nil
}

func (this *SimpleAuth) checkJwt(jwt string) (string, error) {
	Id, err := utils.ParseJwt(jwt, SECRET)

	if err != nil {
		this.Logger.Infof("Token Unabled to Parse: '%s'", err.Error())
		return "", err
	}

	this.Logger.Debugf("Token Parsed with ID: '%s'", Id)

	return Id, nil
}

func (this *SimpleAuth) IsAuthenticated(jwt string) error {
	if _, err := this.checkJwt(jwt); err != nil {
		return err
	}

	return nil
}

func (this *SimpleAuth) IsAuthorized(jwt string, page string) error {
	_, err := this.checkJwt(jwt)
	if err != nil {
		return err
	}

	// check page permission

	return nil
}

func (this SimpleAuth) Login(body string) (*model.LoginResponse, error) {
	var loginRequest model.LoginRequest
	err := this.Validator.MarshalAndValidateLoginRequest(body, loginRequest)
	if err != nil {
		return nil, err
	}

	user, err := this.UserRepo.GetByEmail(loginRequest.Username)

	if !utils.ComparePassword(user.HashedPassword, loginRequest.Password) {
		this.Logger.Infof("INCORRECT PASSWORD ATTEMPT FOR USER '%s' WITH PASSWORD '%s'", loginRequest.Username, loginRequest.Password)
		return nil, types.NewUnauthorizedError()
	}

	return this.generateLoginResponseFromUser(*user)

}

func (this *SimpleAuth) Logout(jwt string) error {
	userId, err := utils.GetIssuerFromJwt(jwt, SECRET)
	if err != nil {
		this.Logger.Errorf("UNSUCCESUFUL LOGOUT WITH TOKEN '%s'", jwt)
	}
	this.Logger.Debugf("SUCCESUFUL LOGOUT FOR USER '%d' TOKEN '%s'", userId, jwt)
	return nil
}

func (this *SimpleAuth) Register(body string) (*model.LoginResponse, error) {
	var loginRequest model.UserRequest
	err := this.Validator.MarshalAndValidateLoginRequest(body, loginRequest)
	if err != nil {
		return nil, err
	}

	if loginRequest.ConfirmPassword != loginRequest.Password {
		return nil, types.NewInvalidInputError()
	}

	user, err := this.UserRepo.Register(loginRequest.FirstName, loginRequest.LastName, loginRequest.Email, loginRequest.Password, loginRequest.RoleID)
	if err != nil {
		return nil, err
	}

	return this.generateLoginResponseFromUser(*user)
}

func (this *SimpleAuth) User(jwt string) (*model.UserResponse, error) {
	issuer, err := utils.GetIssuerFromJwt(jwt, SECRET)
	if err != nil {
		this.Logger.Errorf("CANT RETRIEVE ISSUER FROM TOKEN '%s'", issuer)
		return nil, types.NewInternalServerError()
	}
	userId, err := strconv.ParseUint(issuer, 10, 64)
	if err != nil {
		this.Logger.Errorf("CANT RETRIEVE USERID FOR ISSUER '%s'", issuer)
		return nil, types.NewInternalServerError()
	}
	user, err := this.UserRepo.GetByID(userId)
	if err != nil {
		this.Logger.Errorf("CANT RETRIEVE USER FROM TOKEN '%s'", jwt)
		return nil, types.NewInternalServerError()
	}
	return user, nil
}

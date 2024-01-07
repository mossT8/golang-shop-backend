package repository

import (
	"database/sql"
	"time"

	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/model"
	"tannar.moss/backend/internal/repository/flows"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

type UserRepository interface {
	GetByID(userId uint64) (*model.UserResponse, error)
	GetByEmail(email string) (*model.UserResponse, error)
	Register(firstName string, lastName string, email string, password string, roleId uint64) (*model.UserResponse, error)
	Update(userId uint64, firstName string, lastName string, updatingUserId uint64) (*model.UserResponse, error)
	ResetPassword(userId uint64, newPassword string, updatingUserId uint64) (*model.UserResponse, error)
	ResetEmail(userId uint64, newEmail string, updatingUserId uint64) (*model.UserResponse, error)
	Shutdown()
}

type MySqlUserRepository struct {
	DB     mysql.DbConnection
	Logger logger.Logger
}

const MySystemAutoID = 1

func (this *MySqlUserRepository) Shutdown() {
	err := this.DB.Close()
	if err != nil {
		this.Logger.Errorf("Unabled to close user repo: %s", err.Error())
	}
}

func NewMySqlUserRepository(logger logger.Logger, db mysql.DbConnection) UserRepository {
	return &MySqlUserRepository{
		Logger: logger,
		DB:     db,
	}
}

func (this *MySqlUserRepository) mapStatementToUser(row *sql.Row) (*model.UserResponse, error) {
	var user model.UserResponse
	if row.Err() != nil {
		return nil, types.NewInternalServerError()
	}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword, &user.RoleID, &user.CreatedUser, &user.CreatedAt, &user.UpdatedUser, &user.UpdatedAt, &user.DeletedUser, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			this.Logger.Debugf("No result back for user: %s", err.Error())
			return nil, types.NewNoTFoundOrNoRecordError()
		}
		this.Logger.Errorf("Unabled to marshal user response: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (this *MySqlUserRepository) GetByEmail(email string) (*model.UserResponse, error) {
	query := "SELECT * FROM users WHERE email = ?"
	stmt, err := flows.GetReaderStatement("GetByEmail", query, this.DB, this.Logger)
	if err != nil {
		return nil, err
	}
	this.Logger.Debugf("Running query '%s' with parameter '%s'", query, email)

	result := stmt.QueryRow(email)
	user, err := this.mapStatementToUser(result)

	if err != nil {
		utils.LogExecutingError("GetByEmail", this.Logger, err)
		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (this *MySqlUserRepository) GetByID(userId uint64) (*model.UserResponse, error) {
	query := "SELECT * FROM users WHERE id = ?"
	stmt, err := flows.GetReaderStatement("GetByID", query, this.DB, this.Logger)
	if err != nil {
		return nil, err
	}
	this.Logger.Debugf("Running query '%s' with parameter '%d'", query, userId)
	user, err := this.mapStatementToUser(stmt.QueryRow(userId))
	if err != nil {
		utils.LogExecutingError("GetByID", this.Logger, err)
		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (this *MySqlUserRepository) Register(firstName string, lastName string, email string, password string, roleId uint64) (*model.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		this.Logger.Errorf("Error creating hashPassword: %s", err.Error())
		return nil, types.NewInternalServerError()
	}
	insertedAt := utils.GetCurrentDateFormatedForInsertingIntoDB(time.Now())
	query := "INSERT INTO users (first_name, last_name, email, hashed_password, role_id, created_user, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	this.Logger.Debugf("Running query '%s' with parameter '%s', '%s', '%s', '%v', '%d', '%d' and '%s'", query, firstName, lastName, email, hashedPassword, roleId, MySystemAutoID, insertedAt)
	lastInsertedId, err := flows.PerformEdit(
		"Register",
		query,
		this.DB,
		this.Logger,
		firstName, lastName, email, hashedPassword, roleId, MySystemAutoID, insertedAt)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:          uint64(lastInsertedId),
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		RoleID:      roleId,
		CreatedUser: MySystemAutoID,
		CreatedAt:   insertedAt,
	}, nil
}

func (this *MySqlUserRepository) Update(userId uint64, firstName string, lastName string, updatingUserId uint64) (*model.UserResponse, error) {
	query := "UPDATE users SET first_name = ?, last_name = ?, updated_user = ?, updated_at = now() WHERE id = ?"
	this.Logger.Debugf("Running query '%s' with parameter '%s', '%s', '%d' and '%d'", firstName, lastName, updatingUserId, userId)
	_, err := flows.PerformEdit(
		"Update",
		query,
		this.DB,
		this.Logger,
		firstName, lastName, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return this.GetByID(userId)
}

func (this *MySqlUserRepository) ResetPassword(userId uint64, newPassword string, updatingUserId uint64) (*model.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		this.Logger.Errorf("Error hashing password: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	query := "UPDATE users SET hashed_password = ?, updated_user = ? WHERE id = ?"
	this.Logger.Debugf("Running query '%s' with parameter '%v', '%d' and '%d'", hashedPassword, updatingUserId, userId)
	_, err = flows.PerformEdit(
		"ResetPassword",
		query,
		this.DB,
		this.Logger,
		hashedPassword, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return this.GetByID(userId)
}

func (this *MySqlUserRepository) ResetEmail(userId uint64, newEmail string, updatingUserId uint64) (*model.UserResponse, error) {
	query := "UPDATE users SET email = ?, updated_user = ? WHERE id = ?"
	this.Logger.Debugf("Running query '%s' with parameter '%s', '%d' and '%d'", newEmail, updatingUserId, userId)
	_, err := flows.PerformEdit(
		"ResetEmail",
		query,
		this.DB,
		this.Logger,
		newEmail, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return this.GetByID(userId)
}

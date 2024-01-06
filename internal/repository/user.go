package repository

import (
	"database/sql"

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

func mapStatementToUser(row *sql.Row) (*model.UserResponse, error) {
	var user *model.UserResponse
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.RoleID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (this *MySqlUserRepository) GetByEmail(email string) (*model.UserResponse, error) {
	stmt, err := flows.GetReaderStatement("GetByEmail", "SELECT * FROM users WHERE email = ?", this.DB, this.Logger)
	if err != nil {
		return nil, err
	}

	user, err := mapStatementToUser(stmt.QueryRow(email))
	if err != nil {
		utils.LogExecutingError("GetByEmail", this.Logger, err)
		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (this *MySqlUserRepository) GetByID(userId uint64) (*model.UserResponse, error) {
	stmt, err := flows.GetReaderStatement("GetByID", "SELECT * FROM users WHERE id = ?", this.DB, this.Logger)
	if err != nil {
		return nil, err
	}

	user, err := mapStatementToUser(stmt.QueryRow(userId))
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

	err = flows.PerformEdit(
		"Register",
		"INSERT INTO users (first_name, last_name, email, hashed_password, role_id) VALUES (?, ?, ?, ?, ?)",
		this.DB,
		this.Logger,
		firstName, lastName, email, hashedPassword, roleId)
	if err != nil {
		return nil, err
	}

	return this.GetByEmail(email)
}

func (this *MySqlUserRepository) Update(userId uint64, firstName string, lastName string, updatingUserId uint64) (*model.UserResponse, error) {
	err := flows.PerformEdit(
		"Update",
		"UPDATE users SET first_name = ?, last_name = ?, updated_user = ?, updated_at = now() WHERE id = ?",
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

	err = flows.PerformEdit(
		"ResetPassword",
		"UPDATE users SET hashed_password = ?, updated_user = ? WHERE id = ?",
		this.DB,
		this.Logger,
		hashedPassword, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return this.GetByID(userId)
}

func (this *MySqlUserRepository) ResetEmail(userId uint64, newEmail string, updatingUserId uint64) (*model.UserResponse, error) {
	err := flows.PerformEdit(
		"ResetEmail",
		"UPDATE users SET email = ?, updated_user = ? WHERE id = ?",
		this.DB,
		this.Logger,
		newEmail, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return this.GetByID(userId)
}

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

func (repo *MySqlUserRepository) Shutdown() {
	err := repo.DB.Close()
	if err != nil {
		repo.Logger.Errorf("Unabled to close user repo: %s", err.Error())
	}
}

func NewMySqlUserRepository(logger logger.Logger, db mysql.DbConnection) UserRepository {
	return &MySqlUserRepository{
		Logger: logger,
		DB:     db,
	}
}

func (repo *MySqlUserRepository) mapStatementToUser(row *sql.Row) (*model.UserResponse, error) {
	var user model.UserResponse
	if row.Err() != nil {
		return nil, types.NewInternalServerError()
	}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword, &user.RoleID, &user.CreatedUser, &user.CreatedAt, &user.UpdatedUser, &user.UpdatedAt, &user.DeletedUser, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.Logger.Debugf("No result back for user: %s", err.Error())
			return nil, types.NewNoTFoundOrNoRecordError()
		}
		repo.Logger.Errorf("Unabled to marshal user response: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (repo *MySqlUserRepository) GetByEmail(email string) (*model.UserResponse, error) {
	query := "SELECT * FROM users WHERE email = ?"
	stmt, err := flows.GetReaderStatement("GetByEmail", query, repo.DB, repo.Logger)
	if err != nil {
		return nil, err
	}
	repo.Logger.Debugf("Running query '%s' with parameter '%s'", query, email)

	result := stmt.QueryRow(email)
	user, err := repo.mapStatementToUser(result)

	if err != nil {
		utils.LogExecutingError("GetByEmail", repo.Logger, err)
		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (repo *MySqlUserRepository) GetByID(userId uint64) (*model.UserResponse, error) {
	query := "SELECT * FROM users WHERE id = ?"
	stmt, err := flows.GetReaderStatement("GetByID", query, repo.DB, repo.Logger)
	if err != nil {
		return nil, err
	}
	repo.Logger.Debugf("Running query '%s' with parameter '%d'", query, userId)
	user, err := repo.mapStatementToUser(stmt.QueryRow(userId))
	if err != nil {
		utils.LogExecutingError("GetByID", repo.Logger, err)
		return nil, types.NewInternalServerError()
	}

	return user, nil
}

func (repo *MySqlUserRepository) Register(firstName string, lastName string, email string, password string, roleId uint64) (*model.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		repo.Logger.Errorf("Error creating hashPassword: %s", err.Error())
		return nil, types.NewInternalServerError()
	}
	insertedAt := utils.GetCurrentDateFormatedForInsertingIntoDB(time.Now())
	query := "INSERT INTO users (first_name, last_name, email, hashed_password, role_id, created_user, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	repo.Logger.Debugf("Running query '%s' with parameter '%s', '%s', '%s', '%v', '%d', '%d' and '%s'", query, firstName, lastName, email, hashedPassword, roleId, MySystemAutoID, insertedAt)
	lastInsertedId, err := flows.PerformEdit(
		"Register",
		query,
		repo.DB,
		repo.Logger,
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

func (repo *MySqlUserRepository) Update(userId uint64, firstName string, lastName string, updatingUserId uint64) (*model.UserResponse, error) {
	query := "UPDATE users SET first_name = ?, last_name = ?, updated_user = ?, updated_at = now() WHERE id = ?"
	repo.Logger.Debugf("Running query '%s' with parameter '%s', '%s', '%d' and '%d'", query, firstName, lastName, updatingUserId, userId)
	_, err := flows.PerformEdit(
		"Update",
		query,
		repo.DB,
		repo.Logger,
		firstName, lastName, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(userId)
}

func (repo *MySqlUserRepository) ResetPassword(userId uint64, newPassword string, updatingUserId uint64) (*model.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		repo.Logger.Errorf("Error hashing password: %s", err.Error())
		return nil, types.NewInternalServerError()
	}

	query := "UPDATE users SET hashed_password = ?, updated_user = ? WHERE id = ?"
	repo.Logger.Debugf("Running query '%s' with parameter '%v', '%d' and '%d'", query, hashedPassword, updatingUserId, userId)
	_, err = flows.PerformEdit(
		"ResetPassword",
		query,
		repo.DB,
		repo.Logger,
		hashedPassword, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(userId)
}

func (repo *MySqlUserRepository) ResetEmail(userId uint64, newEmail string, updatingUserId uint64) (*model.UserResponse, error) {
	query := "UPDATE users SET email = ?, updated_user = ? WHERE id = ?"
	repo.Logger.Debugf("Running query '%s' with parameter '%s', '%d' and '%d'", query, newEmail, updatingUserId, userId)
	_, err := flows.PerformEdit(
		"ResetEmail",
		query,
		repo.DB,
		repo.Logger,
		newEmail, updatingUserId, userId)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(userId)
}

package flows

import (
	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

func PerformEdit(queryName string, query string, conn mysql.DbConnection, logger logger.Logger, args ...any) (int64, error) {
	tx, err := conn.GetWriter().Begin()
	if err != nil {
		utils.LogPreparingError(queryName, logger, err)
		return -1, types.NewInternalServerError()
	}

	preparedStmt, err := tx.Prepare(query)
	if err != nil {
		utils.LogPreparingError(queryName, logger, err)
		tx.Rollback()
		return -1, types.NewInternalServerError()
	}
	defer preparedStmt.Close()

	result, err := preparedStmt.Exec(args...)
	if err != nil {
		utils.LogExecutingError(queryName, logger, err)
		tx.Rollback()
		return -1, types.NewInternalServerError()
	}

	err = tx.Commit()
	if err != nil {
		utils.LogCommitError(queryName, logger, err)
		return -1, types.NewInternalServerError()
	}

	id, _ := result.LastInsertId()

	return id, nil
}

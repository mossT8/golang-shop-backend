package flows

import (
	"database/sql"

	"tannar.moss/backend/internal/logger"
	"tannar.moss/backend/internal/repository/mysql"
	"tannar.moss/backend/internal/types"
	"tannar.moss/backend/internal/utils"
)

func GetReaderStatement(queryName string, query string, conn mysql.DbConnection, logger logger.Logger) (*sql.Stmt, error) {
	stmt, err := conn.GetReader().Prepare(query)
	if err != nil {
		utils.LogPreparingError(queryName, logger, err)
		return nil, types.NewInternalServerError()
	}

	return stmt, nil
}

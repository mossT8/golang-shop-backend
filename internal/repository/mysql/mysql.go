package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	Host              string
	Port              int
	RequestTimeout    int
	ConnectionTimeout int
	Dialect           string
	Database          string
	Username          string
	Password          string
}

type MySql interface {
	Ping() error
	Close() error
	GetReader() error
	GetWriter() error
	createConnection(dbConfig DatabaseConfig) (*sql.DB, error)
}

type DbConnection struct {
	readerDB *sql.DB
	writerDB *sql.DB
}

func NewDbConnection(writerCreds, readerCreds DatabaseConfig) (*DbConnection, error) {
	writerDB, err := createConnection(writerCreds)
	if err != nil {
		return nil, err
	}
	readerDB, err := createConnection(readerCreds)
	if err != nil {
		return nil, err
	}
	return &DbConnection{readerDB: readerDB, writerDB: writerDB}, nil
}

func createConnection(dbConfig DatabaseConfig) (*sql.DB, error) {

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	db, err := sql.Open(dbConfig.Dialect, connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(3)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DbConnection) GetReader() *sql.DB {
	return db.readerDB
}

func (db *DbConnection) GetWriter() *sql.DB {
	return db.writerDB
}

func (db *DbConnection) Close() error {
	var err = db.readerDB.Close()
	if err != nil {
		return err
	}

	return db.writerDB.Close()
}

func (db *DbConnection) Ping() error {
	err := db.readerDB.Ping()
	if err != nil {
		return err
	}
	return db.writerDB.Ping()
}

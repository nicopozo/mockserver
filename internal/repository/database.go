package repository

import (
	"database/sql"
	"fmt"
	"time"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	noRowsMessage      = "sql: no rows in result set"
	maxIdleConnections = 100
	maxOpenConnections = 350
	maxLifeTime        = 100 * time.Millisecond
)

//go:generate mockgen -source=client_db.go -destination=../utils/test/mocks/client_db_mock.go -package=mocks
//go:generate mockgen -destination=../utils/test/mocks/sql_result_mock.go -package=mocks database/sql Result

//nolint:gochecknoglobals
var db *sqlx.DB

const defaultForwarderDB = "root:password@/mockserver"

type Database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Prepare(query string) (*sql.Stmt, error)
	Beginx() (*sqlx.Tx, error)
}

func GetDB() (*sqlx.DB, error) {
	var err error
	if db == nil {
		db, err = sqlx.Open("mysql", getDBString()+"?parseTime=true&charset=utf8")

		if err != nil {
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")

			return nil, fmt.Errorf("error connecting DB: %w", err)
		}

		db.SetMaxIdleConns(maxIdleConnections)
		db.SetMaxOpenConns(maxOpenConnections)
		db.SetConnMaxLifetime(maxLifeTime)

		err = db.Ping()
		if err != nil {
			fmt.Printf("########## DB ERROR: " + err.Error() + " #############")

			db = nil
			err = fmt.Errorf("error DB ping: %w", err)
		}
	}

	return db, err
}

func getDBString() string {
	return defaultForwarderDB
}

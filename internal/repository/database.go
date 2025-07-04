package repository

import (
	"database/sql"
	"fmt"
	"os"
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

//go:generate mockgen -destination=../utils/test/mocks/sql_result_mock.go -package=mocks database/sql Result

//nolint:gochecknoglobals
var database *sqlx.DB

type Database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Prepare(query string) (*sql.Stmt, error)
	Beginx() (*sqlx.Tx, error)
}

func GetDB() (*sqlx.DB, error) {
	var err error
	if database == nil {
		database, err = sqlx.Open("mysql", getDBString()+"?parseTime=true&charset=utf8")

		if err != nil {
			fmt.Printf("########## DB ERROR: %s #############\n", err.Error()) //nolint:forbidigo

			return nil, fmt.Errorf("error connecting DB: %w", err)
		}

		database.SetMaxIdleConns(maxIdleConnections)
		database.SetMaxOpenConns(maxOpenConnections)
		database.SetConnMaxLifetime(maxLifeTime)

		for i := 0; i < 10 && (i == 0 || err != nil); i++ {
			if i > 0 {
				time.Sleep(10 * time.Second) //nolint:gomnd
			}

			fmt.Printf("########## CONNECTING TO DB - try i:%v #############\n", i+1) //nolint:forbidigo

			err = database.Ping()
		}

		if err != nil {
			fmt.Printf("########## DB ERROR: %s #############\n", err.Error()) //nolint:forbidigo

			database = nil
			err = fmt.Errorf("error DB ping: %w", err)
		}
	}

	return database, err
}

func getDBString() string {
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/mockserver", user, password, host, port)

	return connStr
}

func getEnv(name, defaultValue string) string {
	if e := os.Getenv(name); e != "" {
		return e
	}

	return defaultValue
}

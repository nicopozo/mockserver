package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	noRowsMessage      = "sql: no rows in result set"
	maxIdleConnections = 100
	maxOpenConnections = 350
	maxLifeTime        = 100 * time.Millisecond

	datasourceMySQL    = "mysql"
	datasourcePostgres = "postgres"
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
	DriverName() string
}

// getDatasource returns "postgres" or "mysql" based on MOCKS_DATASOURCE env var.
func getDatasource() string {
	ds := strings.ToLower(os.Getenv("MOCKS_DATASOURCE"))
	if ds == datasourcePostgres {
		return datasourcePostgres
	}

	return datasourceMySQL
}

func GetDB() (*sqlx.DB, error) {
	if database != nil {
		return database, nil
	}

	return connect()
}

func connect() (*sqlx.DB, error) {
	var err error

	datasource := getDatasource()
	connStr := getMySQLDBString() + "?parseTime=true&charset=utf8"

	if datasource == datasourcePostgres {
		connStr = getPostgresDBString()
	}

	database, err = sqlx.Open(datasource, connStr)
	if err != nil {
		fmt.Printf("########## DB ERROR: %s #############\n", err.Error()) //nolint:forbidigo

		return nil, fmt.Errorf("error connecting DB: %w", err)
	}

	database.SetMaxIdleConns(maxIdleConnections)
	database.SetMaxOpenConns(maxOpenConnections)
	database.SetConnMaxLifetime(maxLifeTime)

	for i := 0; i < 10 && (i == 0 || err != nil); i++ {
		if i > 0 {
			time.Sleep(10 * time.Second) //nolint:mnd
		}

		fmt.Printf("########## CONNECTING TO DB (%s) - try i:%v #############\n", datasource, i+1) //nolint:forbidigo

		err = database.PingContext(context.Background())
	}

	if err != nil {
		fmt.Printf("########## DB ERROR: %s #############\n", err.Error()) //nolint:forbidigo

		database = nil

		return nil, fmt.Errorf("error DB ping: %w", err)
	}

	return database, nil
}

func getMySQLDBString() string {
	if mysqlURL := os.Getenv("MYSQL_URL"); mysqlURL != "" {
		u, err := url.Parse(mysqlURL)
		if err == nil && u.Scheme == "mysql" {
			password, _ := u.User.Password()
			dbName := strings.TrimPrefix(u.Path, "/")

			return fmt.Sprintf("%s:%s@tcp(%s)/%s", u.User.Username(), password, u.Host, dbName)
		}

		return mysqlURL
	}

	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "mockserver")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
}

func getPostgresDBString() string {
	if pgURL := os.Getenv("POSTGRES_URL"); pgURL != "" {
		return pgURL
	}

	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "mockserver")
	sslMode := getEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=mockserver sslmode=%s",
		host, port, user, password, dbName, sslMode)
}

func getEnv(name, defaultValue string) string {
	if e := os.Getenv(name); e != "" {
		return e
	}

	return defaultValue
}

// FormatQuery translates a MySQL-style query to the target driver dialect.
// For postgres: replaces backtick-escaped identifiers with double quotes and
// rebinds positional parameters from ? to $1, $2, etc.
func FormatQuery(query string, driver string) string {
	if driver != datasourcePostgres {
		return query
	}

	// Replace backtick-wrapped identifiers with double-quoted ones
	re := regexp.MustCompile("`([^`]+)`")
	query = re.ReplaceAllString(query, `"$1"`)

	// Rebind ? placeholders to $1, $2, ...
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	return query
}

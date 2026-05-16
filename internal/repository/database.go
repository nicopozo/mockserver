package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicopozo/mockserver/internal/configs"
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

type Database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Prepare(query string) (*sql.Stmt, error)
	Beginx() (*sqlx.Tx, error)
	DriverName() string
}

func NewSQLDB(cfg *configs.Config) (*sqlx.DB, error) {
	return connect(cfg)
}

func connect(cfg *configs.Config) (*sqlx.DB, error) {
	var err error

	datasource := strings.ToLower(cfg.DataSource)
	connStr := getMySQLDBString(cfg) + "?parseTime=true&charset=utf8"

	if datasource == datasourcePostgres {
		connStr = getPostgresDBString(cfg)
	}

	databaseConn, err := sqlx.Open(datasource, connStr)
	if err != nil {
		fmt.Printf("########## DB ERROR: %s #############\n", err.Error()) //nolint:forbidigo

		return nil, fmt.Errorf("error connecting DB: %w", err)
	}

	databaseConn.SetMaxIdleConns(maxIdleConnections)
	databaseConn.SetMaxOpenConns(maxOpenConnections)
	databaseConn.SetConnMaxLifetime(maxLifeTime)

	for i := 0; i < 10 && (i == 0 || err != nil); i++ {
		if i > 0 {
			time.Sleep(10 * time.Second) //nolint:mnd
		}

		fmt.Printf("########## CONNECTING TO DB (%s) - try i:%v #############\n", datasource, i+1) //nolint:forbidigo

		err = databaseConn.PingContext(context.Background())
	}

	if err != nil {
		fmt.Printf("########## DB ERROR: %s #############\n", err.Error()) //nolint:forbidigo

		return nil, fmt.Errorf("error DB ping: %w", err)
	}

	return databaseConn, nil
}

func getMySQLDBString(cfg *configs.Config) string {
	if cfg.Database.URL != "" && strings.HasPrefix(cfg.Database.URL, "mysql://") {
		u, err := url.Parse(cfg.Database.URL)
		if err == nil {
			password, _ := u.User.Password()
			dbName := strings.TrimPrefix(u.Path, "/")

			return fmt.Sprintf("%s:%s@tcp(%s)/%s", u.User.Username(), password, u.Host, dbName)
		}
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
}

func getPostgresDBString(cfg *configs.Config) string {
	if cfg.Database.URL != "" && !strings.HasPrefix(cfg.Database.URL, "mysql://") {
		return cfg.Database.URL
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=mockserver sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
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

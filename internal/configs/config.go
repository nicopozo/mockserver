package configs

import (
	"os"
	"strings"
)

type Config struct {
	DataSource string
	MocksFile  string
	Database   DatabaseConfig
	Dynamo     DynamoConfig
	AWS        AWSConfig
}

type DatabaseConfig struct {
	URL      string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

type DynamoConfig struct {
	Endpoint    string
	TablePrefix string
}

type AWSConfig struct {
	Region string
}

func New() *Config {
	return &Config{
		DataSource: getEnv("MOCKS_DATASOURCE", "file"),
		MocksFile:  getEnv("MOCKS_FILE", "/tmp/mocks.json"),
		Database: DatabaseConfig{
			URL:      getEnv("MYSQL_URL", getEnv("POSTGRES_URL", "")),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     getEnv("DB_NAME", "mockserver"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Dynamo: DynamoConfig{
			Endpoint:    os.Getenv("DYNAMO_ENDPOINT"),
			TablePrefix: getEnv("DYNAMO_TABLE_PREFIX", "mockserver_"),
		},
		AWS: AWSConfig{
			Region: getEnv("AWS_REGION", "us-east-1"),
		},
	}
}

func getEnv(name, defaultValue string) string {
	if e := os.Getenv(name); e != "" {
		return e
	}

	return defaultValue
}

func (c *Config) IsSQL() bool {
	ds := strings.ToLower(c.DataSource)

	return ds == "mysql" || ds == "postgres"
}

func (c *Config) IsDynamo() bool {
	ds := strings.ToLower(c.DataSource)

	return ds == "dynamodb" || ds == "dynamo"
}

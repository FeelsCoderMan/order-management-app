package config

import (
	"fmt"
	"os"
	"strings"
)

type PostgresConfig struct {
	User     string
	Password string
	Db       string
	Host     string
	Port     string
}

func LoadPostgresConfig() (PostgresConfig, error) {
	var missing []string

	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		missing = append(missing, "POSTGRES_USER")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		missing = append(missing, "POSTGRES_PASSWORD")
	}

	db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		missing = append(missing, "POSTGRES_DB")
	}

	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		missing = append(missing, "POSTGRES_HOST")
	}

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		missing = append(missing, "POSTGRES_PORT")
	}

	if len(missing) > 0 {
		return PostgresConfig{}, fmt.Errorf(
			"Missing required environment variables: %s",
			strings.Join(missing, ","),
		)
	}

	return PostgresConfig{
		User:     user,
		Password: password,
		Db:       db,
		Host:     host,
		Port:     port,
	}, nil
}

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AuthConfig struct {
	AccessSecret   string
	AccessExpires  time.Duration
}

func LoadAuthConfig() (AuthConfig, error) {
	var missing []string

	secret, ok := os.LookupEnv("JWT_ACCESS_SECRET")
	if !ok {
		missing = append(missing, "JWT_ACCESS_SECRET")
	}

	accessExpiresString, ok := os.LookupEnv("JWT_ACCESS_EXPIRES_IN_MINUTES")
	if !ok {
		missing = append(missing, "JWT_ACCESS_EXPIRES_IN_MINUTES")
	}


	if len(missing) > 0 {
		return AuthConfig{}, fmt.Errorf(
			"Missing required environment variables: %s",
			strings.Join(missing, ","),
		)
	}

	accessExpires, err := strconv.Atoi(accessExpiresString)
	if err != nil {
		return AuthConfig{}, fmt.Errorf("JWT_ACCESS_EXPIRED_IN_MINUTES must be an integer: %w", err)
	}

	return AuthConfig{
		AccessSecret: secret,
		AccessExpires: time.Duration(accessExpires) * time.Minute,
	}, nil
}

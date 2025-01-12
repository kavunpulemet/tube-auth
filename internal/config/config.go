package config

import (
	"fmt"
	"os"
	"time"
)

const (
	serverPortEnv = "SERVER_PORT"
	dbHostEnv     = "DB_HOST"
	dbPortEnv     = "DB_PORT"
	dbUserEnv     = "DB_USER"
	dbNameEnv     = "DB_NAME"
	dbPasswordEnv = "DB_PASSWORD"
	dbSSLModeEnv  = "DB_SSLMODE"
)

type Config struct {
	ServerPort         string
	DBConnectionString string
	JWTSecret          []byte
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

func NewConfig() (Config, error) {
	accessTTLStr := os.Getenv("ACCESS_TOKEN_TTL")
	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil || accessTTL <= 0 {
		return Config{}, fmt.Errorf("invalid ACCESS_TOKEN_TTL: %s", accessTTLStr)
	}

	refreshTTLStr := os.Getenv("REFRESH_TOKEN_TTL")
	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil || refreshTTL <= 0 {
		return Config{}, fmt.Errorf("invalid REFRESH_TOKEN_TTL: %s", refreshTTLStr)
	}

	return Config{
		ServerPort: fmt.Sprintf(":%s", os.Getenv(serverPortEnv)),
		DBConnectionString: fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			os.Getenv(dbHostEnv), os.Getenv(dbPortEnv), os.Getenv(dbUserEnv),
			os.Getenv(dbNameEnv), os.Getenv(dbPasswordEnv), os.Getenv(dbSSLModeEnv)),
		JWTSecret:       []byte(os.Getenv("JWT_SECRET")),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}, nil
}

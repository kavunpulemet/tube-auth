package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	type env struct {
		serverPort      string
		dbHost          string
		dbPort          string
		dbUser          string
		dbName          string
		dbPassword      string
		dbSSLMode       string
		jwtSecret       string
		accessTokenTTL  string
		refreshTokenTTL string
	}

	setEnv := func(env env) {
		os.Setenv("SERVER_PORT", env.serverPort)
		os.Setenv("DB_HOST", env.dbHost)
		os.Setenv("DB_PORT", env.dbPort)
		os.Setenv("DB_USER", env.dbUser)
		os.Setenv("DB_NAME", env.dbName)
		os.Setenv("DB_PASSWORD", env.dbPassword)
		os.Setenv("DB_SSLMODE", env.dbSSLMode)
		os.Setenv("JWT_SECRET", env.jwtSecret)
		os.Setenv("ACCESS_TOKEN_TTL", env.accessTokenTTL)
		os.Setenv("REFRESH_TOKEN_TTL", env.refreshTokenTTL)
	}

	tests := []struct {
		name    string
		env     env
		want    Config
		wantErr bool
	}{
		{
			name: "valid config",
			env: env{
				serverPort:      "80",
				dbHost:          "localhost",
				dbPort:          "5432",
				dbUser:          "test_user",
				dbName:          "test_db",
				dbPassword:      "password",
				dbSSLMode:       "disable",
				jwtSecret:       "secret",
				accessTokenTTL:  "15m",
				refreshTokenTTL: "720h",
			},

			want: Config{
				ServerPort:         ":80",
				DBConnectionString: "host=localhost port=5432 user=test_user dbname=test_db password=password sslmode=disable",
				JWTSecret:          []byte("secret"),
				AccessTokenTTL:     15 * time.Minute,
				RefreshTokenTTL:    720 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "invalid ACCESS_TOKEN_TTL",
			env: env{
				serverPort:      "80",
				dbHost:          "localhost",
				dbPort:          "5432",
				dbUser:          "test_user",
				dbName:          "test_db",
				dbPassword:      "password",
				dbSSLMode:       "disable",
				jwtSecret:       "secret",
				accessTokenTTL:  "invalid",
				refreshTokenTTL: "720h",
			},

			want:    Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.env)

			got, err := NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

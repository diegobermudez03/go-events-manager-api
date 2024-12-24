package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     	string
	DbConfig 	*DbConfig
	AuthConfig 	*AuthConfig
	EmailConfig *EmailConfig
}

type AuthConfig struct{
	SecondsLife 			int64
	AccessTokenExpiration 	int64
	JWTSecret				string
}

type DbConfig struct {
	Addr         string
	MaxOpenConn  int
	MaxIdleConns int
	MaxIdleTime  string
}

type EmailConfig struct{
	ApiKey 		string
}


func NewConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", ":8081"),
		DbConfig: NewDBConfig(),
		AuthConfig: NewAuthConfig(),
		EmailConfig: NewEmailConfig(),
	}
}

func NewDBConfig() *DbConfig {
	return &DbConfig{
		Addr: getEnv("POSTGRES_URL", "postgres://admin:secret@localhost:5432/events_go?sslmode=disable"),
	}
}

func NewAuthConfig() *AuthConfig{
	return &AuthConfig{
		SecondsLife: getEnvAsInt("REFRESH_TOKEN_LIFE_HOURS", 1440),
		AccessTokenExpiration: getEnvAsInt("ACCESS_TOKEN_LIFE_SECONDS", 600),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
}

func NewEmailConfig() *EmailConfig{
	return &EmailConfig{
		ApiKey: getEnv("EMAIL_API_KEY", ""),
	}
}


func getEnv(param string, fallback string) string {
	if val, ok := os.LookupEnv(param); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(param string, fallback int64) int64 {
	if val, ok := os.LookupEnv(param); ok {
		number, err := strconv.Atoi(val)
		if err != nil{
			return fallback
		}
		return int64(number)
	}
	return fallback
}
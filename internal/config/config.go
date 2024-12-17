package config

import "os"

type Config struct {
	Port     string
	DbConfig *DbConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConn  int
	MaxIdleConns int
	MaxIdleTime  string
}

func NewConfig() *Config {
	return &Config{
		Port:     getEnv("PORT", ":8081"),
		DbConfig: NewDBConfig(),
	}
}

func NewDBConfig() *DbConfig {
	return &DbConfig{
		Addr: getEnv("POSTGRES_URL", "postgres://admin:secret@localhost:5432/events_go?sslmode=disable"),
	}
}

func getEnv(param string, fallback string) string {
	if val, ok := os.LookupEnv(param); ok {
		return val
	}
	return fallback
}
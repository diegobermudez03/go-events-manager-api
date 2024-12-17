package config

type Config struct {
	Port     string
	DbConfig DbConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConn  int
	MaxIdleConns int
	MaxIdleTime  string
}
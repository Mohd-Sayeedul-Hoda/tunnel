package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type Config struct {
	Server struct {
		Port int
		Host string
	}
	DB struct {
		DSN          string
		MaxOpenConn  int
		MaxIdealConn int
		MaxIdleTime  string
	}
	Cache struct {
		DSN string
	}
	Token struct {
		AccessTokenPublicKey   string
		AccessTokenPrivateKey  string
		AccessTokenExpiredIn   string
		AccessTokenMaxAge      uint
		RefreshTokenPublicKey  string
		RefreshTokenPrivateKey string
		RefreshTokenExpiredIn  string
		RefreshTokenMaxAge     uint
	}
	AppVersion    int    // app version like 1
	AppEnviroment string //  production|development|debug
	Debug         bool   // run code in debug mode mostly debug log will be displayed
}

func (c *Config) validate() error {
	if c.DB.DSN == "" {
		return errors.New("DB_DSN is not set")
	}
	if c.Cache.DSN == "" {
		return errors.New("REDIS_DSN is not set")
	}
	if c.Token.AccessTokenPublicKey == "" {
		return errors.New("ACCESS_TOKEN_PUBLIC_KEY is not set")
	}
	if c.Token.AccessTokenPrivateKey == "" {
		return errors.New("ACCESS_TOKEN_PRIVATE_KEY is not set")
	}
	if c.Token.RefreshTokenPublicKey == "" {
		return errors.New("REFRESH_TOKEN_PUBLIC_KEY is not set")
	}
	if c.Token.RefreshTokenPrivateKey == "" {
		return errors.New("REFRESH_TOKEN_PRIVATE_KEY is not set")
	}
	return nil
}

func InitializeConfig(getenv func(string) string, args []string) (*Config, error) {
	cfg := Config{}

	cfg.Server.Port = getEnvInt(getenv, "PORT", 8000)
	cfg.Server.Host = getEnvString(getenv, "HOST", "localhost")

	cfg.DB.DSN = getEnvString(getenv, "DB_DSN", "")
	cfg.DB.MaxOpenConn = getEnvInt(getenv, "DB-MAX-OPEN-CONNS", 10)
	cfg.DB.MaxIdealConn = getEnvInt(getenv, "DB-MAX-IDLE-CONNS", 10)
	cfg.DB.MaxIdleTime = getEnvString(getenv, "DB-MAX-IDLE-TIME", "10m")

	cfg.Cache.DSN = getEnvString(getenv, "REDIS_DSN", "")

	cfg.AppVersion = getEnvInt(getenv, "APP_VERSION", 1)
	cfg.AppEnviroment = getEnvString(getenv, "APP_ENVIROMENT", "development")

	cfg.Token.AccessTokenPublicKey = getEnvString(getenv, "ACCESS_TOKEN_PUBLIC_KEY", "")
	cfg.Token.AccessTokenPrivateKey = getEnvString(getenv, "ACCESS_TOKEN_PRIVATE_KEY", "")
	cfg.Token.AccessTokenExpiredIn = getEnvString(getenv, "ACCESS_TOKEN_EXPIRED_IN", "")
	cfg.Token.AccessTokenMaxAge = uint(getEnvInt(getenv, "ACCESS_TOKEN_MAXAGE", 15))

	cfg.Token.RefreshTokenPublicKey = getEnvString(getenv, "REFRESH_TOKEN_PUBLIC_KEY", "")
	cfg.Token.RefreshTokenPrivateKey = getEnvString(getenv, "REFRESH_TOKEN_PRIVATE_KEY", "")
	cfg.Token.RefreshTokenExpiredIn = getEnvString(getenv, "REFRESH_TOKEN_EXPIRED_IN", "")
	cfg.Token.RefreshTokenMaxAge = uint(getEnvInt(getenv, "REFRESH_TOKEN_MAXAGE", 60))

	flag.BoolVar(&cfg.Debug, "debug", false, "Debug mode")

	err := flag.CommandLine.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.AppEnviroment = "debug"
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getEnvString(getenv func(string) string, key string, fallback string) string {
	value := getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(getenv func(string) string, key string, fallback int) int {
	value := getenv(key)
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}

func getEnvBool(getenv func(string) string, key string, fallback bool) bool {
	value := getenv(key)
	if value == "" {
		return fallback
	}

	lowerValue := strings.ToLower(value)
	if lowerValue == "true" || lowerValue == "1" || lowerValue == "yes" || lowerValue == "y" {
		return true
	}
	if lowerValue == "false" || lowerValue == "0" || lowerValue == "no" || lowerValue == "n" {
		return false
	}

	return fallback
}

func getEnvSlice(getenv func(string) string, key string, fallback []string) []string {
	value := getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}
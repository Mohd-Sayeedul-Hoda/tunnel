package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
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
		AccessTokenExpiredIn   time.Duration
		AccessTokenMaxAge      uint
		RefreshTokenPublicKey  string
		RefreshTokenPrivateKey string
		RefreshTokenExpiredIn  time.Duration
		RefreshTokenMaxAge     uint
	}
	AppVersion        int           // app version like 1
	AppEnv            string        //  prod|dev|debug
	Debug             bool          // run code in debug mode mostly debug log will be displayed
	EmailOtpExpiredIn time.Duration // after how much time email expired token get expired
	EmailOtpSalt      string
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
	if c.EmailOtpSalt == "" {
		return errors.New("EMAIL_OTP_SALT is not set")
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
	cfg.AppEnv = getEnvString(getenv, "APP_ENV", "development")

	cfg.Token.AccessTokenPublicKey = getEnvString(getenv, "ACCESS_TOKEN_PUBLIC_KEY", "")
	cfg.Token.AccessTokenPrivateKey = getEnvString(getenv, "ACCESS_TOKEN_PRIVATE_KEY", "")
	cfg.Token.AccessTokenMaxAge = uint(getEnvInt(getenv, "ACCESS_TOKEN_MAXAGE", 15))

	cfg.Token.RefreshTokenPublicKey = getEnvString(getenv, "REFRESH_TOKEN_PUBLIC_KEY", "")
	cfg.Token.RefreshTokenPrivateKey = getEnvString(getenv, "REFRESH_TOKEN_PRIVATE_KEY", "")
	cfg.Token.RefreshTokenMaxAge = uint(getEnvInt(getenv, "REFRESH_TOKEN_MAXAGE", 60))

	accessExpireIn := getEnvString(getenv, "ACCESS_TOKEN_EXPIRED_IN", "15m")
	refreshExpireIn := getEnvString(getenv, "REFRESH_TOKEN_EXPIRED_IN", "60m")

	accessTokenExpireIn, err := time.ParseDuration(accessExpireIn)
	if err != nil {
		return nil, fmt.Errorf("invalid access token duration: %w", err)
	}

	refreshTokenExpireIn, err := time.ParseDuration(refreshExpireIn)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token duration: %w", err)
	}

	cfg.Token.AccessTokenExpiredIn = accessTokenExpireIn
	cfg.Token.RefreshTokenExpiredIn = refreshTokenExpireIn

	cfg.EmailOtpSalt = getEnvString(getenv, "EMAIL_OTP_SALT", "")
	emailOtpExpiredIn := getEnvString(getenv, "EMAIL_OTP_EXPIRED_IN", "15m")

	cfg.EmailOtpExpiredIn, err = time.ParseDuration(emailOtpExpiredIn)
	if err != nil {
		return nil, fmt.Errorf("invalid email otp expiration duration: %w", err)
	}

	flag.BoolVar(&cfg.Debug, "debug", false, "Debug mode")

	err = flag.CommandLine.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		cfg.AppEnv = "debug"
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

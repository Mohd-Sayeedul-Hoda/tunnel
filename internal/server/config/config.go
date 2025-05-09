package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
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
	}
	AppVersion    int
	AppEnviroment string
	Debug         bool
}

func InitalizeConfig(getenv func(string) string, args []string) (*Config, error) {
	config := Config{}
	flag.IntVar(&config.Server.Port, "port", 8000, "port for http server")
	flag.StringVar(&config.Server.Host, "host", "localhost", "host for http server")
	flag.IntVar(&config.DB.MaxOpenConn, "db-max-open-conns", 15, "PostgreSQL max open connection")
	flag.IntVar(&config.DB.MaxIdealConn, "db-max-idle-conns", 15, "PostgreSQL max idle connections")
	flag.StringVar(&config.AppEnviroment, "app-env", "Development", "production|development|debug")
	flag.BoolVar(&config.Debug, "debug", false, "Debug mode")
	err := flag.CommandLine.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	if config.Debug {
		config.AppEnviroment = "debug"
	}

	config.DB.DSN = getenv("DB_DSN")
	if config.DB.DSN == "" {
		return nil, fmt.Errorf("no database connection string found")
	}
	appVersion := getenv("APP_VERSION")
	config.AppVersion, err = strconv.Atoi(appVersion)
	if err != nil {
		return nil, errors.New("app version should be integer")
	}

	return &config, nil
}

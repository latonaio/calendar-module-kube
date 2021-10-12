package app

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type DatabaseEnv struct {
	Addr     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"3306"`
	User     string `envconfig:"DB_USER" default:"latona"`
	Password string `envconfig:"DB_PASSWORD"`
	Name     string `envconfig:"DB_NAME" default:"calendar"`
}

type ServerEnv struct {
	Port string `envconfig:"SERVER_PORT" default:"localhost"`
	Addr string `envconfig:"SERVER_HOST" default:"8888"`
}

func GetDatabaseEnv() (*DatabaseEnv, error) {
	var env DatabaseEnv
	if err := envconfig.Process("DB", &env); err != nil {
		return nil, fmt.Errorf("can not get database config\n")
	}

	return &env, nil
}

func GetServerEnv() (*ServerEnv, error) {
	var env ServerEnv
	if err := envconfig.Process("SERVER", &env); err != nil {
		return nil, fmt.Errorf("can not get server envconfig\n")
	}

	return &env, nil
}

package env

import (
	env "github.com/Netflix/go-env"
)

type Environment struct {
	Database struct {
		User     string `env:"DATABASE_USER"`
		Host     string `env:"DATABASE_HOST"`
		Port     string `env:"DATABASE_PORT"`
		Password string `env:"DATABASE_PASSWORD"`
		Dbname   string `env:"DATABASE_DBNAME"`
		Sslmode  string `env:"DATABASE_SSLMODE"`
	}
}

func NewEnvironment() (*Environment, error) {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

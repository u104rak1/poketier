package env

import (
	"fmt"

	envpkg "github.com/caarlos0/env/v11"
)

type Env struct {
	APP_PORT string `env:"APP_PORT" envDefault:"8080"`

	POSTGRES_HOST     string `env:"POSTGRES_HOST" envDefault:"postgres"`
	POSTGRES_DBNAME   string `env:"POSTGRES_DBNAME" envDefault:"POCGO_LOCAL_DB"`
	POSTGRES_USER     string `env:"POSTGRES_USER" envDefault:"local_user"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT" envDefault:"5432"`
	POSTGRES_SSLMODE  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`

	LOG_LEVEL     string `env:"LOG_LEVEL" envDefault:"debug"`
	IS_SILENT_LOG bool   `env:"IS_SILENT_LOG" envDefault:"false"`
}

func NewEnv() *Env {
	e, err := envpkg.ParseAs[Env]()
	if err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}
	return &e
}

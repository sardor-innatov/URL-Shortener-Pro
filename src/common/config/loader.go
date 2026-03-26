package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func Load[T any](instance *T, filenames ...string) error {
	if len(filenames) > 0 {
		if err := godotenv.Load(filenames...); err != nil {
			return err
		}
	}

	if err := env.Parse(instance); err != nil {
		return err
	}

	return nil
}

func MustLoad(instance *EnvProject, filenames ...string) {
	if err := Load(instance); err != nil {
		panic(err)
	}
}
package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Env            string
		SentryEndpoint string
		PgConfig
	}

	PgConfig struct {
		PgUser          string
		PgPassword      string
		PgDatabase      string
		PostgresPort    string
		PostgresHost    string
		PostgresSSLMode string
	}
)

func Load(env string) *Config {
	if env == "" {
		env = "local"
	}

	viper.AutomaticEnv()
	viper.SetConfigName(fmt.Sprintf("config_%s", env))
	viper.SetConfigType("json")
	viper.AddConfigPath("./../config") // test
	viper.AddConfigPath("./config")    // local
	viper.AddConfigPath("/app/config") // inside container

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &Config{
		Env:            env,
		SentryEndpoint: viper.GetString("SENTRY_ENDPOINT"),
		PgConfig: PgConfig{
			PgUser:          viper.GetString("POSTGRES_USER"),
			PgPassword:      viper.GetString("POSTGRES_PASSWORD"),
			PgDatabase:      viper.GetString("POSTGRES_DB"),
			PostgresPort:    viper.GetString("POSTGRES_PORT"),
			PostgresHost:    viper.GetString("POSTGRES_HOST"),
			PostgresSSLMode: viper.GetString("POSTGRES_SSL"),
		},
	}
}

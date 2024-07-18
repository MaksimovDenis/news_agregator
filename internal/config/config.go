package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	APIPort      int    `default:"8080"`
	PgConnString string `default:"postgres://admin:admin@localhost:5432/newsAgregator?sslmode=disable"`
	LogLevel     string `default:"debug"`
}

// InitConfig init config
func InitConfig() (*Config, error) {
	var cnf Config
	err := envconfig.Process("SKILLFACTORY", &cnf)
	if err != nil {
		return nil, err
	}

	return &cnf, nil
}

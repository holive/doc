package doc

import (
	"os"

	"github.com/holive/doc/app/config"
)

func loadConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if connectionString != "" {
		cfg.Mongo.URI = connectionString
	}

	return cfg, nil
}

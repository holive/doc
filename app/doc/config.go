package doc

import "github.com/holive/doc/app/config"

func loadConfig() (*config.Config, error) {
	return config.New()
}

package doc

import (
	"github.com/holive/doc/app/config"
	"github.com/holive/doc/app/docApi"
	"github.com/holive/doc/app/mongo"
)

func initMongoClient(cfg *config.Config) (*mongo.Client, error) {
	return mongo.New(&mongo.ClientConfig{
		URI:      cfg.Mongo.URI,
		Database: cfg.Mongo.Database,
		Timeout:  cfg.Mongo.Timeout,
	})
}

func initDocApiService(client *mongo.Client) *docApi.Service {
	repository := mongo.NewDocApiRepository(client)

	return docApi.NewService(repository)
}

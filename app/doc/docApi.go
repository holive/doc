package doc

import (
	"github.com/holive/doc/app/docApi"
	"github.com/holive/doc/app/mongo"
)

func initDocApiService(client *mongo.Client) *docApi.Service {
	repository := mongo.NewDocApiRepository(client)

	return docApi.NewService(repository)
}

package doc

import (
	"github.com/holive/doc/app/mongo"
	"github.com/holive/doc/app/squads"
)

func initSquadsService(client *mongo.Client) *squads.Service {
	repository := mongo.NewSquadsRepository(client)

	return squads.NewService(repository)
}

package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/holive/doc/app/squads"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type SquadsRepository struct {
	collection *mongo.Collection
}

func (sr *SquadsRepository) Create(ctx context.Context, squad squads.Squad) (squads.Squad, error) {
	var exists squads.Squad
	_ = sr.collection.FindOne(ctx, bson.M{"name": squad.Name}).Decode(&exists)

	if exists.Name == squad.Name {
		return squads.Squad{}, errors.New("squad already exists")
	}

	resp, err := sr.collection.InsertOne(ctx, squad)
	if err != nil {
		return squads.Squad{}, errors.Wrap(err, "could not create the squad")
	}

	var s squads.Squad
	if err = sr.collection.FindOne(ctx, bson.M{"_id": resp.InsertedID}).Decode(&s); err != nil {
		return squads.Squad{}, errors.Wrap(err, "could not find the new squad")
	}

	return s, nil
}

func (sr *SquadsRepository) GetByKey(ctx context.Context, key string) (squads.Squad, error) {
	var s squads.Squad

	filter := bson.M{
		"key": bson.M{"$eq": key},
	}

	if err := sr.collection.FindOne(ctx, filter).Decode(&s); err != nil {
		return squads.Squad{}, errors.Wrap(err, "could not find the created squad")
	}

	return s, nil
}

func NewSquadsRepository(conn *Client) *SquadsRepository {
	return &SquadsRepository{
		collection: conn.db.Collection(SquadsCollection),
	}
}

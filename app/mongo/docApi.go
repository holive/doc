package mongo

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/holive/doc/app/docApi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocApiRepository struct {
	collection *mongo.Collection
}

func (dr *DocApiRepository) Create(ctx context.Context, doc *docApi.DocApi) error {
	_, err := dr.collection.InsertOne(ctx, doc)
	if err != nil {
		return errors.Wrap(err, "could not create a feed")
	}

	return nil
}

func (dr *DocApiRepository) Find(ctx context.Context, squad string, projeto string, versao string) (*docApi.DocApi, error) {
	var f docApi.DocApi

	filter := bson.M{
		"squad":   bson.M{"$eq": squad},
		"projeto": bson.M{"$eq": projeto},
		"versao":  bson.M{"$eq": versao},
	}

	if err := dr.collection.FindOne(ctx, filter).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func (dr *DocApiRepository) Delete(ctx context.Context, squad string, projeto string, versao string) error {
	filter := bson.M{
		"squad":   bson.M{"$eq": squad},
		"projeto": bson.M{"$eq": projeto},
		"versao":  bson.M{"$eq": versao},
	}

	_, err := dr.collection.DeleteOne(ctx, filter)

	return err
}

func (dr *DocApiRepository) FindAll(ctx context.Context, limit string, offset string) (*docApi.SearchResult, error) {
	intLimit, intOffset, err := dr.getLimitOffset(limit, offset)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset)

	cur, err := dr.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return &docApi.SearchResult{}, err
	}

	total, err := dr.collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := dr.resultFromCursor(ctx, cur)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return &docApi.SearchResult{
		Docs: results,
		Result: struct {
			Offset int64 `json:"offset"`
			Limit  int64 `json:"limit"`
			Total  int64 `json:"total"`
		}{
			Offset: intOffset,
			Limit:  intLimit,
			Total:  total,
		},
	}, nil
}

func (dr *DocApiRepository) getLimitOffset(limit string, offset string) (int64, int64, error) {
	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "24"
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return 0, 0, err
	}

	return int64(intLimit), int64(intOffset), nil
}

func (dr *DocApiRepository) resultFromCursor(ctx context.Context, cur *mongo.Cursor) ([]docApi.DocApi, error) {
	var results []docApi.DocApi
	for cur.Next(ctx) {
		var elem docApi.DocApi
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	return results, nil
}

func NewDocApiRepository(conn *Client) *DocApiRepository {
	return &DocApiRepository{
		collection: conn.db.Collection(DocApiCollection),
	}
}

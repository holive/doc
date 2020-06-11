package mongo

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/holive/doc/app/docApi"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocApiRepository struct {
	collection *mongo.Collection
}

const (
	Squad   = "squad"
	Projeto = "projeto"
	Versao  = "versao"
)

func (dr *DocApiRepository) Create(ctx context.Context, doc *docApi.DocApi) error {
	filter := bson.M{
		"squad":   bson.M{"$eq": doc.Squad},
		"projeto": bson.M{"$eq": doc.Projeto},
		"versao":  bson.M{"$eq": doc.Versao},
	}

	bsonDoc, err := bson.Marshal(doc)
	if err != nil {
		return errors.Wrap(err, "could not marshal bson")
	}

	opts := options.Replace().SetUpsert(true)

	_, err = dr.collection.ReplaceOne(ctx, filter, bsonDoc, opts)
	if err != nil {
		return errors.Wrap(err, "could not create a doc")
	}

	return nil
}

func (dr *DocApiRepository) Find(ctx context.Context, squad string, projeto string, versao string) (*docApi.DocApi, error) {
	var f docApi.DocApi

	filter := bson.M{
		Squad:   bson.M{"$eq": squad},
		Projeto: bson.M{"$eq": projeto},
		Versao:  bson.M{"$eq": versao},
	}

	if err := dr.collection.FindOne(ctx, filter).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func (dr *DocApiRepository) Delete(ctx context.Context, squad string, projeto string, versao string) error {
	filter := bson.M{
		Squad:   bson.M{"$eq": squad},
		Projeto: bson.M{"$eq": projeto},
		Versao:  bson.M{"$eq": versao},
	}

	_, err := dr.collection.DeleteOne(ctx, filter)

	return err
}

func (dr *DocApiRepository) FindAll(ctx context.Context, limit string, offset string) (*docApi.SearchResult, error) {
	intLimit, intOffset, err := dr.getLimitOffset(limit, offset)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset).SetProjection(bson.M{
		Squad:   1,
		Projeto: 1,
		Versao:  1,
	})

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

	return dr.returnSearchResult(results, intOffset, intLimit, total), nil
}

func (dr *DocApiRepository) FindBySquad(ctx context.Context, squad string, limit string, offset string) (*docApi.SearchResult, error) {
	intLimit, intOffset, err := dr.getLimitOffset(limit, offset)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset).SetProjection(bson.M{
		Squad:   1,
		Projeto: 1,
		Versao:  1,
	})

	filter := bson.M{
		Squad: bson.M{"$eq": squad},
	}

	cur, err := dr.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return &docApi.SearchResult{}, err
	}

	total, err := dr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := dr.resultFromCursor(ctx, cur)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return dr.returnSearchResult(results, intOffset, intLimit, total), nil
}

func (dr *DocApiRepository) SearchProject(ctx context.Context, project string, limit string, offset string) (*docApi.SearchResult, error) {
	intLimit, intOffset, err := dr.getLimitOffset(limit, offset)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset).SetProjection(bson.M{
		Squad:   1,
		Projeto: 1,
		Versao:  1,
	})

	filter := bson.M{
		Projeto: primitive.Regex{Pattern: project, Options: "i"},
	}

	cur, err := dr.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return &docApi.SearchResult{}, err
	}

	total, err := dr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := dr.resultFromCursor(ctx, cur)
	if err != nil {
		return &docApi.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return dr.returnSearchResult(results, intOffset, intLimit, total), nil
}

func (dr *DocApiRepository) getLimitOffset(limit string, offset string) (int64, int64, error) {
	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "6"
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

func (dr *DocApiRepository) returnSearchResult(results []docApi.DocApi, offset int64, limit int64, total int64) *docApi.SearchResult {
	return &docApi.SearchResult{
		Docs: results,
		Result: struct {
			Offset int64 `json:"offset"`
			Limit  int64 `json:"limit"`
			Total  int64 `json:"total"`
		}{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
	}
}

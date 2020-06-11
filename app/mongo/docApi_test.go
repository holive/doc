package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedRepository_Create(t *testing.T) {
	client, err := NewConnectionTest()
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	docApirepo := NewDocApiRepository(NewDBTest(client))
	err = docApirepo.collection.Drop(context.Background())
	require.NoError(t, err)

	{
		r, err := docApirepo.collection.InsertOne(
			context.Background(),
			bson.D{
				{Squad, "squadXpto"},
				{Projeto, "alfred"},
				{Versao, "v1"},
			})
		require.NoError(t, err)
		require.NotNil(t, r.InsertedID)

		r, err = docApirepo.collection.InsertOne(
			context.Background(),
			bson.D{
				{Squad, "squadXpto"},
				{Projeto, "sherlock"},
				{Versao, "v1"},
			})
		require.NoError(t, err)
		require.NotNil(t, r.InsertedID)
	}

	result, err := docApirepo.SearchProject(context.Background(), "she", "", "")
	require.NoError(t, err)
	assert.True(t, result.Result.Total != 0)
	assert.Equal(t, "sherlock", result.Docs[0].Projeto)

	result, err = docApirepo.SearchProject(context.Background(), "alf", "", "")
	require.NoError(t, err)
	assert.Equal(t, "alfred", result.Docs[0].Projeto)
}

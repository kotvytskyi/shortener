package mongo

import (
	"context"
	"testing"

	"github.com/kotvytskyi/testmongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreate(t *testing.T) {
	coll, teardown := testmongo.CreateTestMongoConnection(t)
	defer teardown()

	repo := Short{coll}

	err := repo.Create(context.Background(), "test_url", "test_short")
	require.Nil(t, err)

	dto := &shortUrlDto{}
	coll.FindOne(context.Background(), bson.M{"url": "test_url"}).Decode(dto)

	require.Equal(t, "test_url", dto.Url)
	require.Equal(t, "test_short", dto.Short)
}

func TestGetUrl(t *testing.T) {
	coll, teardown := testmongo.CreateTestMongoConnection(t)
	defer teardown()

	coll.InsertOne(context.Background(), shortUrlDto{Url: "test", Short: "t"})

	repo := Short{coll}
	got, err := repo.GetUrl(context.Background(), "t")
	require.Nil(t, err)

	require.Equal(t, "test", got)
}

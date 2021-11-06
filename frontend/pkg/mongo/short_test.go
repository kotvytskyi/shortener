package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/kotvytskyi/testmongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	coll, teardown := testmongo.CreateTestMongoConnection(ctx, t)
	defer teardown()

	repo := Short{coll}

	err := repo.Create(ctx, "test_url", "test_short")
	require.Nil(t, err)

	dto := &shortURLDto{}
	coll.FindOne(ctx, bson.M{"url": "test_url"}).Decode(dto)

	require.Equal(t, "test_url", dto.URL)
	require.Equal(t, "test_short", dto.Short)
}

func TestGetUrl(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	coll, teardown := testmongo.CreateTestMongoConnection(ctx, t)
	defer teardown()

	coll.InsertOne(ctx, shortURLDto{URL: "test", Short: "t"})

	repo := Short{coll}
	got, err := repo.GetURL(ctx, "t")
	require.Nil(t, err)

	require.Equal(t, "test", got)
}

package mongo

import (
	"context"
	"testing"

	"github.com/kotvytskyi/shortener/testutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreate(t *testing.T) {
	coll, teardown := testutils.CreateTestMongoConnection(t)
	defer teardown()

	repo := Short{coll}

	err := repo.Create(context.Background(), "test_url", "test_short")
	assert.Nil(t, err)

	dto := &shortUrlDto{}
	coll.FindOne(context.Background(), bson.M{"url": "test_url"}).Decode(dto)

	assert.Equal(t, "test_url", dto.Url)
	assert.Equal(t, "test_short", dto.Short)
}

func TestGetUrl(t *testing.T) {
	coll, teardown := testutils.CreateTestMongoConnection(t)
	defer teardown()

	coll.InsertOne(context.Background(), shortUrlDto{Url: "test", Short: "t"})

	repo := Short{coll}
	got, err := repo.GetUrl(context.Background(), "t")
	assert.Nil(t, err)

	assert.Equal(t, "test", got)
}

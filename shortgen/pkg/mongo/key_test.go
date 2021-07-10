package mongo

import (
	"context"
	"testing"

	"github.com/kotvytskyi/shortener/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	coll, teardown := testutils.CreateTestMongoConnection(t)
	defer teardown()

	mongo := &KeyRepository{
		coll: coll,
	}

	want := Key{
		Value: "test_key",
	}
	err := mongo.Create(context.Background(), want.Value)
	assert.Nil(t, err)

	got := Key{}
	coll.FindOne(context.Background(), want).Decode(&got)

	assert.Equal(t, got, want)
}

package generator

import (
	"context"
	"testing"

	"github.com/kotvytskyi/shortener/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	coll, teardown := testutils.CreateTestMongoConnection(t)
	defer teardown()

	mongo := &Mongo{
		coll: coll,
	}

	want := Key{
		Value: "test_key",
	}
	err := mongo.create(context.Background(), want)
	assert.Nil(t, err)

	got := Key{}
	coll.FindOne(context.Background(), want).Decode(&got)

	assert.Equal(t, got, want)
}

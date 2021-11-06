package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/kotvytskyi/testmongo"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	coll, teardown := testmongo.CreateTestMongoConnection(ctx, t)
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

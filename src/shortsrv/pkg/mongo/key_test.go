package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/kotvytskyi/testmongo"
	"github.com/stretchr/testify/assert"
)

func TestReserveKey(t *testing.T) {
	t.Run("returns err when no keys in pool", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		coll, teardown := testmongo.CreateTestMongoConnection(ctx, t)
		defer teardown()

		mongo := &KeyRepository{coll}
		_, err := mongo.ReserveKey(context.Background())
		assert.NotNil(t, err)
	})

	t.Run("returns err when no unused keys in pool", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		coll, teardown := testmongo.CreateTestMongoConnection(ctx, t)
		defer teardown()

		mongo := &KeyRepository{coll}
		mongo.coll.InsertOne(context.Background(), &keyDto{Value: "test", Used: true})

		_, err := mongo.ReserveKey(context.Background())
		assert.NotNil(t, err)
	})

	t.Run("returns a key", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		coll, teardown := testmongo.CreateTestMongoConnection(ctx, t)
		defer teardown()

		mongo := &KeyRepository{coll}
		key := &keyDto{Value: "test", Used: false}
		mongo.coll.InsertOne(context.Background(), key)

		reserved, err := mongo.ReserveKey(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, key.Value, reserved)
	})
}

package app

import (
	"context"
	"testing"

	"github.com/kotvytskyi/shortener/testutils"
	"github.com/stretchr/testify/assert"
)

func TestReserveKey(t *testing.T) {
	t.Run("returns err when no keys in pool", func(t *testing.T) {
		coll, teardown := testutils.CreateTestMongoConnection(t)
		defer teardown()

		mongo := &Mongo{coll}
		_, err := mongo.ReserveKey(context.Background())
		assert.NotNil(t, err)
	})

	t.Run("returns err when no unused keys in pool", func(t *testing.T) {
		coll, teardown := testutils.CreateTestMongoConnection(t)
		defer teardown()

		mongo := &Mongo{coll}
		mongo.coll.InsertOne(context.Background(), &Key{Value: "test", Used: true})

		_, err := mongo.ReserveKey(context.Background())
		assert.NotNil(t, err)
	})

	t.Run("returns err when no unused keys in pool", func(t *testing.T) {
		coll, teardown := testutils.CreateTestMongoConnection(t)
		defer teardown()

		mongo := &Mongo{coll}
		key := &Key{Value: "test", Used: false}
		mongo.coll.InsertOne(context.Background(), key)

		reserved, err := mongo.ReserveKey(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, key, reserved)
	})
}

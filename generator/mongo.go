package generator

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	coll *mongo.Collection
}

type MongoConfig struct {
	Endpoint string
}

func NewMongo(config MongoConfig) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Endpoint))
	if err != nil {
		return nil, err
	}

	mongo := &Mongo{}
	mongo.coll = client.Database("shortener").Collection("keys_pool")
	return mongo, nil
}

func (r *Mongo) create(ctx context.Context, key Key) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	_, err := r.coll.InsertOne(timeoutCtx, key)

	return err
}

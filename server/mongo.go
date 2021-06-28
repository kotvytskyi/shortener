package server

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (m *Mongo) ReserveKey(ctx context.Context) (*Key, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	keys := m.coll
	filter := bson.M{"$or": bson.A{
		bson.M{"used": false},
		bson.M{"used": bson.M{"$exists": false}}}}

	update := bson.M{
		"$set": bson.M{"used": true},
	}

	key := keys.FindOneAndUpdate(ctx, filter, update)
	if key.Err() != nil {
		return nil, key.Err()
	}

	result := &Key{}
	key.Decode(&result)
	return result, nil
}

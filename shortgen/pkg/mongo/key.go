package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Key struct {
	Value string `json:"val" bson:"val"`
}

type Mongo struct {
	coll *mongo.Collection
}

type Config struct {
	Username string
	Password string
	Address  string
}

func NewKeyRepository(config Config) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	endpoint := fmt.Sprintf("mongodb://%s:%s@%s:27017", config.Username, config.Password, config.Address)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		return nil, err
	}

	mongo := &Mongo{}
	mongo.coll = client.Database("shortener").Collection("keys_pool")
	return mongo, nil
}

func (r *Mongo) Create(ctx context.Context, key string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	dto := Key{Value: key}
	_, err := r.coll.InsertOne(timeoutCtx, dto)

	return err
}

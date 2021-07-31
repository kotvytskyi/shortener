package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type keyDto struct {
	Value string `json:"val" bson:"val"`
	Used  bool   `json:"-" bson:"used"`
}

type KeyRepository struct {
	coll *mongo.Collection
}

type Config struct {
	Address  string
	User     string
	Password string
}

const ConnectionTimeout = 10 * time.Second
const Timeout = 2 * time.Second

func NewKeyRepository(config Config) (*KeyRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	endpoint := fmt.Sprintf("mongodb://%s:%s@%s:27017", config.User, config.Password, config.Address)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		return nil, err
	}

	mongo := &KeyRepository{}
	mongo.coll = client.Database("shortener").Collection("keys_pool")

	return mongo, nil
}

func (m *KeyRepository) ReserveKey(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
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
		return "", key.Err()
	}

	result := &keyDto{}

	err := key.Decode(&result)
	if err != nil {
		return "", key.Err()
	}

	return result.Value, nil
}

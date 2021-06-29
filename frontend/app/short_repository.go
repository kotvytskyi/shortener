package app

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoShortRepository struct {
	coll *mongo.Collection
}

type MongoParams struct {
	Endpoint string
}

type shortUrlDto struct {
	Url   string `bson:"url"`
	Short string `bson:"short"`
}

func NewMongoShortRepository(params MongoParams) (*MongoShortRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(params.Endpoint))
	if err != nil {
		return nil, err
	}

	coll := client.Database("shortener").Collection("shorts")

	return &MongoShortRepository{coll}, nil
}

func (r *MongoShortRepository) Create(ctx context.Context, url string, short string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	dto := shortUrlDto{url, short}
	_, err := r.coll.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

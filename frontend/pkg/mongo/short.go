package mongo

import (
	"context"
	"time"

	"github.com/kotvytskyi/frontend/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Short struct {
	Coll *mongo.Collection
}

type shortUrlDto struct {
	Url   string `bson:"url"`
	Short string `bson:"short"`
}

func NewShort(cfg config.MongoConfig) (*Short, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	p := NewParams(cfg)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(p.Endpoint))
	if err != nil {
		return nil, err
	}

	coll := client.Database("shortener").Collection("shorts")

	return &Short{coll}, nil
}

func (r *Short) Create(ctx context.Context, url string, short string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	dto := shortUrlDto{url, short}
	_, err := r.Coll.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

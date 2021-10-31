package mongo

import (
	"context"
	"time"

	"github.com/kotvytskyi/frontend/pkg/config"
	"github.com/kotvytskyi/frontend/pkg/shorter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ConnectionTimeout = time.Second * 10
const Timeout = time.Second * 2

type Short struct {
	Coll *mongo.Collection
}

type shortURLDto struct {
	URL   string `bson:"url"`
	Short string `bson:"short"`
}

func NewShort(cfg config.MongoConfig) (*Short, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	p := NewParams(cfg)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(p.Endpoint))
	if err != nil {
		return nil, err
	}

	coll := client.Database("shortener").Collection("shorts")

	return &Short{coll}, nil
}

func (r *Short) GetURL(ctx context.Context, short string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	var dto shortURLDto
	err := r.Coll.FindOne(ctx, bson.M{"short": short}).Decode(&dto)

	if err != nil && err == mongo.ErrNoDocuments {
		return "", &shorter.NotFoundError{}
	}

	return dto.URL, nil
}

func (r *Short) Create(ctx context.Context, url string, short string) error {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	dto := shortURLDto{url, short}

	_, err := r.Coll.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

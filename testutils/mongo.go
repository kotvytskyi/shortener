package testutils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTestMongoConnection(t *testing.T) (coll *mongo.Collection, teardown func()) {
	url := getMongoUrl(t)
	client, err := driver.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		t.Errorf("An error occurred during creation of mongo connection: %v", err)
	}

	db := client.Database("test")
	collName := fmt.Sprintf("test_%d", time.Now().Nanosecond())
	coll = db.Collection(collName)

	teardown = func() {
		coll.Drop(context.Background())
		client.Disconnect(context.Background())
	}

	_ = coll.Drop(context.Background())
	return coll, teardown
}

func getMongoUrl(t *testing.T) string {
	url := os.Getenv("MONGO_TEST")
	if url == "" {
		url = "mongodb://localhost:27017"
		t.Logf("No MONGO_TEST in env, defaulted to %s", url)
	}

	return url
}

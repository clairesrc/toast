package main

// mongodb client
import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoConnectTimeoutSeconds = 5

// mongoClient is a wrapper around the mongo.Client type
// that allows us to mock the mongo.Client type in tests
// by embedding this type in another struct
// and overriding the methods we need to mock
type mongoClient struct {
	*mongo.Client
}

func newMongoClient() (*mongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeoutSeconds*time.Second)
	defer cancel()

	if os.Getenv("MONGODB_URI") == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable not set")
	}

	mongoURI := os.Getenv("MONGODB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("cannot open database connection: %w", err)
	}

	return &mongoClient{client}, nil
}

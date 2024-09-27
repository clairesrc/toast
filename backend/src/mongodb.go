package main

// mongodb client
import (
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoClient is a wrapper around the mongo.Client type
// that allows us to mock the mongo.Client type in tests
// by embedding this type in another struct
// and overriding the methods we need to mock
type mongoClient struct {
	*mongo.Client
}

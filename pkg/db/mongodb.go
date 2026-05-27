// Package db provides reusable database client helpers.
package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoClient wraps the official MongoDB client and exposes the target database.
type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoClient connects to MongoDB at uri, selects dbName, and verifies the
// connection with a Ping. The caller is responsible for calling Disconnect.
func NewMongoClient(ctx context.Context, uri, dbName string) (*MongoClient, error) {
	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("mongodb connect: %w", err)
	}

	if err := client.Ping(connectCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongodb ping: %w", err)
	}

	return &MongoClient{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

// Disconnect cleanly closes the MongoDB connection.
func (m *MongoClient) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents the MongoDB client connection
type MongoDB struct {
	Client    *mongo.Client
	Database  *mongo.Database
	UsersCol  *mongo.Collection
	TokensCol *mongo.Collection
	ResetCol  *mongo.Collection
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	return &MongoDB{
		Client:    client,
		Database:  database,
		UsersCol:  database.Collection("users"),
		TokensCol: database.Collection("blacklisted_tokens"),
		ResetCol:  database.Collection("reset_tokens"),
	}, nil
} // Close closes the MongoDB connection
func (m *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		log.Printf("❌ Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("🔌 Disconnected from MongoDB")
	}
}

package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	*mongo.Client
}

func NewConnection(ctx context.Context, connectionURI string) (*mongoClient, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		return nil, err
	}
	return &mongoClient{client}, nil
}

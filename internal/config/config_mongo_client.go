package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoClient(ctx context.Context, cfg MongoConfig) (*mongo.Client, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI()))
	if err != nil {
		return nil, fmt.Errorf("connect mongodb: %w", err)
	}

	if err := client.Ping(timeoutCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	return client, nil
}

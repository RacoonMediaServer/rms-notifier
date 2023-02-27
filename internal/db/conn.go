package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const databaseTimeout = 40 * time.Second

type Database struct {
	cli      *mongo.Client
	db       *mongo.Database
	settings *mongo.Collection
	journal  *mongo.Collection
}

func Connect(uri string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect to db failed: %w", err)
	}

	if err = cli.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("connect to db failed: %w", err)
	}

	notifier := cli.Database("notifier")

	db := &Database{
		cli:      cli,
		db:       notifier,
		settings: notifier.Collection("settings"),
		journal:  notifier.Collection("journal"),
	}

	return db, nil
}

package db

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d Database) LoadSettings(ctx context.Context) (*rms_notifier.NotifierSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	found := d.settings.FindOne(ctx, bson.D{{}})
	if found.Err() != nil {
		if errors.Is(found.Err(), mongo.ErrNoDocuments) {
			return &config.DefaultSettings, nil
		}
		return nil, found.Err()
	}

	var result rms_notifier.NotifierSettings
	return &result, found.Decode(&result)
}

func (d Database) SaveSettings(ctx context.Context, settings *rms_notifier.NotifierSettings) error {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	_, err := d.settings.ReplaceOne(ctx, bson.D{{}}, settings, opts)
	return err
}

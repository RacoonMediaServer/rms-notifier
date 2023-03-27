package db

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (d Database) StoreEvent(ctx context.Context, sender string, e interface{}) error {
	record := Event{Sender: sender}
	switch e := e.(type) {
	case *events.Notification:
		record.Notification = e
		record.Timestamp = time.Now()
	case *events.Malfunction:
		record.Malfunction = e
		record.Timestamp = time.Unix(e.Timestamp, 0)
	case *events.Alert:
		record.Alert = e
		record.Timestamp = time.Unix(e.Timestamp, 0)
	default:
		return errors.New("unknown Event type")
	}

	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	_, err := d.journal.InsertOne(ctx, &record)
	return err
}

func (d Database) LoadEvents(ctx context.Context, from, to *time.Time, limit uint) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	sort := 1

	filter := bson.D{}
	if from != nil && to != nil {
		if from.Before(*to) {
			filter = bson.D{
				{"$and",
					bson.A{
						bson.D{{"timestamp", bson.D{{"$gt", *from}}}},
						bson.D{{"timestamp", bson.D{{"$lt", *to}}}},
					}},
			}
		} else {
			filter = bson.D{
				{"$and",
					bson.A{
						bson.D{{"timestamp", bson.D{{"$gt", *to}}}},
						bson.D{{"timestamp", bson.D{{"$lt", *from}}}},
					}},
			}
			sort = -1
		}
	} else if from != nil {
		filter = bson.D{{"timestamp", bson.D{{"$gt", *from}}}}
	} else if to != nil {
		filter = bson.D{{"timestamp", bson.D{{"$lt", *to}}}}
	} else {
		sort = -1
	}

	opts := options.Find().SetSort(bson.D{{"timestamp", sort}})
	cur, err := d.journal.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	var results []*Event
	for cur.Next(ctx) {
		e := Event{}
		if err = cur.Decode(&e); err != nil {
			return nil, err
		}
		results = append(results, &e)
		if uint(len(results)) >= limit {
			break
		}
	}
	return results, nil
}

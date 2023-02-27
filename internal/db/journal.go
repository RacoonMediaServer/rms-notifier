package db

import (
	"context"
	"errors"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"time"
)

func (d Database) StoreEvent(ctx context.Context, sender string, e interface{}) error {
	record := event{Sender: sender}
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
		return errors.New("unknown event type")
	}

	ctx, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	_, err := d.journal.InsertOne(ctx, &record)
	return err
}

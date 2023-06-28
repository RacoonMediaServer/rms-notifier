package db

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"time"
)

type Event struct {
	Notification *events.Notification
	Malfunction  *events.Malfunction
	Alert        *events.Alert
	Timestamp    time.Time
}

func (e Event) Unpack() interface{} {
	if e.Notification != nil {
		return e.Notification
	} else if e.Malfunction != nil {
		return e.Malfunction
	} else if e.Alert != nil {
		return e.Alert
	}

	return nil
}

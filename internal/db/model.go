package db

import (
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"time"
)

type event struct {
	rms_notifier.Event
	Sender    string
	Timestamp time.Time
}

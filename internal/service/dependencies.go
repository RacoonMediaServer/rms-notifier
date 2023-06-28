package service

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/db"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-notifier/internal/notifier"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"time"
)

type Database interface {
	LoadSettings(ctx context.Context) (*rms_notifier.NotifierSettings, error)
	SaveSettings(ctx context.Context, settings *rms_notifier.NotifierSettings) error
	StoreEvent(ctx context.Context, e interface{}) error
	LoadEvents(ctx context.Context, from, to *time.Time, limit uint) ([]*db.Event, error)
}

type Formatter interface {
	Format(event interface{}) (*formatter.Message, error)
}

type Notifier interface {
	SetSettings(settings notifier.Settings)
	Notify(topic string, message *formatter.Message)
}

package service

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-notifier/internal/notifier"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
)

type Database interface {
	LoadSettings(ctx context.Context) (*rms_notifier.Settings, error)
	SaveSettings(ctx context.Context, settings *rms_notifier.Settings) error
	StoreEvent(ctx context.Context, sender string, e interface{}) error
}

type Formatter interface {
	Format(sender string, event interface{}) (*formatter.Message, error)
}

type Notifier interface {
	SetSettings(settings notifier.Settings)
	Notify(topic string, message *formatter.Message)
}

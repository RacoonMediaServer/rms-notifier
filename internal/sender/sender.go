package sender

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
)

// Sender represents entity which could send messages to some external systems
type Sender interface {
	Send(ctx context.Context, message *formatter.Message) error
}

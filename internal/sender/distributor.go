package sender

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
)

// Distributor distribute messages by topics
type Distributor struct {
	topics map[string]Sender
}

func (d *Distributor) Add(topic string, sender Sender) {
	if d.topics == nil {
		d.topics = map[string]Sender{}
	}
	d.topics[topic] = sender
}

func (d *Distributor) Send(ctx context.Context, topic string, message *formatter.Message) error {
	sender, ok := d.topics[topic]
	if !ok {
		return nil
	}
	return sender.Send(ctx, message)
}

package sender

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
)

// Composite is a Sender which broadcast messages to a few Sender's
type Composite struct {
	children []Sender
}

func (c *Composite) Add(sender Sender) {
	c.children = append(c.children, sender)
}

func (c *Composite) Send(ctx context.Context, message *formatter.Message) error {
	var compositeErr error
	for _, s := range c.children {
		if err := s.Send(ctx, message); err != nil {
			compositeErr = fmt.Errorf("%s: %w", compositeErr, err)
		}
	}
	return compositeErr
}

package notifier

import (
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"go-micro.dev/v4/logger"
)

func (n *Notifier) process() {
	for {
		select {
		case cmd := <-n.cmd:
			n.processCommand(cmd)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Notifier) processCommand(cmd interface{}) {
	switch content := cmd.(type) {
	case *Settings:
		n.setSettings(content)
	case *formatter.Message:
		n.notify(content)
	default:
		logger.Errorf("Unknown command type: %T", content)
	}
}

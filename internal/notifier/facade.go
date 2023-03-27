package notifier

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-notifier/internal/sender"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"go-micro.dev/v4/logger"
	"sync"
	"time"
)

const notifyTimeout = 2 * time.Minute

// Notifier is entity for deliver notifications to users
type Notifier struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	cmd chan interface{}
	f   sender.Factory

	sender sender.Distributor
}

// Settings are configuration of all notifications
type Settings struct {
	Rules map[string]*rms_notifier.NotifierSettings_Rules
}

func New(f sender.Factory) *Notifier {
	n := Notifier{
		cmd: make(chan interface{}),
		f:   f,
	}
	n.ctx, n.cancel = context.WithCancel(context.Background())

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		go n.process()
	}()

	return &n
}

func (n *Notifier) SetSettings(settings Settings) {
	n.cmd <- &settings
}

func (n *Notifier) Notify(topic string, message *formatter.Message) {
	n.cmd <- &notify{topic: topic, msg: message}
}

func (n *Notifier) Stop() {
	n.cancel()
	n.wg.Wait()
}

func (n *Notifier) setSettings(settings *Settings) {
	n.sender = sender.Distributor{}
	for topic, rules := range settings.Rules {
		composite := sender.Composite{}
		for _, rule := range rules.Rule {
			composite.Add(n.f.New(rule.Method, rule.Destination))
		}
		n.sender.Add(topic, &composite)
	}
}

func (n *Notifier) notify(notify *notify) {
	ctx, cancel := context.WithTimeout(n.ctx, notifyTimeout)
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer cancel()

		if err := n.sender.Send(ctx, notify.topic, notify.msg); err != nil {
			logger.Errorf("Notify failed: %s", err)
		}
	}()
}

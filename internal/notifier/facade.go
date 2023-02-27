package notifier

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
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
	f   servicemgr.ServiceFactory
}

// Settings are configuration of all notifications
type Settings struct {
	TelegramEnabled bool
	Rules           map[string]*rms_notifier.Settings_Rules
}

func New(f servicemgr.ServiceFactory) *Notifier {
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

func (n *Notifier) Notify(message *formatter.Message) {
	n.cmd <- message
}

func (n *Notifier) Stop() {
	n.cancel()
	n.wg.Wait()
}

func (n *Notifier) setSettings(settings *Settings) {
}

func (n *Notifier) notify(content *formatter.Message) {
	ctx, cancel := context.WithTimeout(n.ctx, notifyTimeout)
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer cancel()
		<-ctx.Done()
	}()
}

package sender

import (
	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

// Factory can create Sender's
type Factory interface {
	New(method rms_notifier.Rule_Method, destination string) Sender
}

type factory struct {
	f      servicemgr.ServiceFactory
	remote config.Remote
}

func (f factory) New(method rms_notifier.Rule_Method, destination string) Sender {
	switch method {
	case rms_notifier.Rule_Telegram:
		return newTelegramSender(f.f)
	default:
		panic("unknown notification method")
	}
}

func NewFactory(f servicemgr.ServiceFactory, remote config.Remote) Factory {
	return &factory{
		f:      f,
		remote: remote,
	}
}

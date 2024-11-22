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
	f       servicemgr.ServiceFactory
	remote  config.Remote
	apiKey  string
	backend config.EmailBackend
}

func (f factory) New(method rms_notifier.Rule_Method, destination string) Sender {
	switch method {
	case rms_notifier.Rule_Telegram:
		return newTelegramSender(f.f)
	case rms_notifier.Rule_Email:
		if f.backend == config.EmailBackend_TrueNAS {
			return newTruenasEmailSender(f.remote, f.apiKey, destination)
		}
		return newRmsEmailSender(f.remote, f.apiKey, destination)
	case rms_notifier.Rule_SMS:
		return newSmsSender(f.remote, f.apiKey, destination)
	default:
		panic("unknown notification method")
	}
}

func NewFactory(f servicemgr.ServiceFactory, conf config.Configuration) Factory {
	return &factory{
		f:       f,
		remote:  conf.Remote,
		apiKey:  conf.APIKey,
		backend: conf.EmailBackend,
	}
}

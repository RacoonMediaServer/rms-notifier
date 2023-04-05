package sender

import (
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-post/pkg/client/client"
	"github.com/RacoonMediaServer/rms-post/pkg/client/client/notify"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"unicode/utf8"
)

type smsSender struct {
	auth runtime.ClientAuthInfoWriter
	cli  *client.Client
	tel  string
}

func (s smsSender) Send(ctx context.Context, message *formatter.Message) error {
	const textLimit = 60
	subject := message.Subject
	if utf8.RuneCountInString(message.Subject) > textLimit {
		runes := []rune(subject)
		runes = runes[:textLimit]
		subject = string(runes)
	}

	req := notify.NotifySMSParams{
		Text:    subject,
		To:      s.tel,
		Context: ctx,
	}

	_, err := s.cli.Notify.NotifySMS(&req, s.auth)
	if err != nil {
		return err
	}

	return nil
}

func newSmsSender(remote config.Remote, device string, tel string) Sender {
	s := smsSender{tel: tel}
	tr := httptransport.New(remote.Host, remote.Path, []string{remote.Scheme})
	s.auth = httptransport.APIKeyAuth("X-Token", "header", device)
	s.cli = client.New(tr, strfmt.Default)
	return &s
}

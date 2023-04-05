package sender

import (
	"bytes"
	"context"
	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-post/pkg/client/client"
	"github.com/RacoonMediaServer/rms-post/pkg/client/client/notify"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type emailSender struct {
	auth   runtime.ClientAuthInfoWriter
	cli    *client.Client
	mailTo string
}

func (s emailSender) Send(ctx context.Context, message *formatter.Message) error {
	req := notify.NotifyEmailParams{
		Subject: message.Subject,
		Text:    message.BodyHtml,
		To:      s.mailTo,
		Context: ctx,
	}

	if message.Attachment != nil && message.Attachment.Type == communication.Attachment_Photo {
		req.Attachment = runtime.NamedReader("image_0.jpg", bytes.NewReader(message.Attachment.Content))
	}

	_, err := s.cli.Notify.NotifyEmail(&req, s.auth)
	if err != nil {
		return err
	}

	return nil
}

func newEmailSender(remote config.Remote, device string, mailTo string) Sender {
	s := emailSender{mailTo: mailTo}
	tr := httptransport.New(remote.Host, remote.Path, []string{remote.Scheme})
	s.auth = httptransport.APIKeyAuth("X-Token", "header", device)
	s.cli = client.New(tr, strfmt.Default)
	return &s
}

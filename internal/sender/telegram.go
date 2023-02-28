package sender

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

type telegramSender struct {
	f servicemgr.ServiceFactory
}

func (t telegramSender) Send(ctx context.Context, message *formatter.Message) error {
	msg := communication.BotMessage{}
	msg.Text = fmt.Sprintf("<b>%s</b>\n\n%s", message.Subject, message.BodyPlain)
	msg.Attachment = message.Attachment
	_, err := t.f.NewBotClient().SendMessage(ctx, &rms_bot_client.SendMessageRequest{Message: &msg})
	return err
}

func newTelegramSender(f servicemgr.ServiceFactory) Sender {
	return &telegramSender{f: f}
}

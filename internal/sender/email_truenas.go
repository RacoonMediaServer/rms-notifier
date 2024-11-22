package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"go-micro.dev/v4/logger"
)

type truenasEmailSender struct {
	mailTo string
	apiKey string
	url    string
}

type truenasMailMessage struct {
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
	To      []string `json:"to"`
}

type truenasSendRequest struct {
	MailMessage truenasMailMessage `json:"mail_message"`
}

func (s truenasEmailSender) Send(ctx context.Context, message *formatter.Message) error {
	mail := truenasSendRequest{
		MailMessage: truenasMailMessage{
			Subject: message.Subject,
			Html:    message.BodyHtml,
			To:      []string{s.mailTo},
		},
	}

	body, err := json.Marshal(&mail)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url, bytes.NewReader(body))
	logger.Info(s.url)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	logger.Info("Send")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func newTruenasEmailSender(remote config.Remote, apiKey string, mailTo string) Sender {
	u := fmt.Sprintf("%s://%s:%d/api/v2.0/mail/send", remote.Scheme, remote.Host, remote.Port)
	return &truenasEmailSender{mailTo: mailTo, apiKey: apiKey, url: u}
}

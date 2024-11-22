package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
)

type truenasEmailSender struct {
	mailTo string
	apiKey string
	url    string
}

type truenasMailMessage struct {
	Subject     string   `json:"subject"`
	Html        string   `json:"html"`
	To          []string `json:"to"`
	Attachments bool     `json:"attachments"`
}

type truenasDataPart struct {
	MailMessage truenasMailMessage `json:"mail_message"`
}

type truenasFilePart struct {
	Headers []map[string]interface{} `json:"headers"`
	Content []byte                   `json:"content"`
}

func (s truenasEmailSender) Send(ctx context.Context, message *formatter.Message) error {
	mail := truenasDataPart{
		MailMessage: truenasMailMessage{
			Subject:     message.Subject,
			Html:        message.BodyHtml,
			To:          []string{s.mailTo},
			Attachments: message.Attachment != nil,
		},
	}

	contentType := "application/json"

	body, err := json.Marshal(&mail)
	if err != nil {
		return err
	}

	if message.Attachment != nil {
		fp := []truenasFilePart{
			{
				Headers: []map[string]interface{}{
					{
						"name":  "Content-Transfer-Encoding",
						"value": "base64",
					},
					{
						"name":   "Content-Type",
						"value":  message.Attachment.MimeType,
						"params": map[string]string{"name": "image_0.jpg"},
					},
				},
				Content: message.Attachment.Content,
			},
		}
		file, err := json.Marshal(&fp)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(buf)
		defer writer.Close()

		dataPart, err := writer.CreateFormField("data")
		if err != nil {
			return err
		}
		_, err = dataPart.Write(body)
		if err != nil {
			return err
		}
		filePart, err := writer.CreateFormField("file")
		if err != nil {
			return err
		}
		_, err = filePart.Write(file)
		if err != nil {
			return err
		}
		contentType = writer.FormDataContentType()
		body = buf.Bytes()
	}

	req, err := http.NewRequest("POST", s.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	req.Header.Add("Content-Type", contentType)

	client := &http.Client{}
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

package formatter

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"go-micro.dev/v4/logger"
	"html/template"
	plain "text/template"
	"time"
)

//go:embed templates
var templates embed.FS

var htmlTemplates *template.Template
var textTemplates *plain.Template

func init() {
	htmlTemplates = template.Must(template.ParseFS(templates, "templates/*.html"))
	textTemplates = plain.Must(plain.ParseFS(templates, "templates/*.txt"))
}

type Formatter struct {
}

func prettySender(sender string) string {
	switch sender {
	case "rms-torrent":
		return "Торрент-клиент"
	case "rms-cctv":
		return "Система видеонаблюдения"
	default:
		return sender
	}
}

func (f Formatter) Format(sender string, event interface{}) (*Message, error) {
	switch e := event.(type) {
	case *events.Notification:
		return f.formatNotification(sender, e), nil
	case *events.Malfunction:
		return f.formatMalfunction(sender, e), nil
	case *events.Alert:
		return f.formatAlert(sender, e), nil
	default:
		return nil, fmt.Errorf("unknown event type: %T", e)
	}
}

func (f Formatter) formatNotification(sender string, e *events.Notification) *Message {
	m := &Message{}
	switch e.Kind {
	case events.Notification_DownloadComplete:
		m.Severity = Info
		m.Subject = "Загрузка завершена"
	case events.Notification_DownloadFailed:
		m.Severity = Fault
		m.Subject = "Ошибка при загрузке контента"
	case events.Notification_TranscodingDone:
		m.Severity = Info
		m.Subject = "Транскодирование завершено"
	case events.Notification_TranscodingFailed:
		m.Severity = Fault
		m.Subject = "Транскодирование завершено с ошибкой"
	case events.Notification_TorrentRemoved:
		// системное уведомление, так что не отправляем
		return nil
	default:
		logger.Errorf("Unknown notification code: %s", e.Kind)
		m.Subject = "Уведомление"
	}

	ctx := uiContext{
		Title:  m.Subject,
		Sender: prettySender(sender),
		Kind:   e.Kind.String(),
	}

	if e.ItemTitle != nil {
		ctx.Item = *e.ItemTitle
	}

	// TODO: обработка истории с транскодированием видео

	var buf bytes.Buffer
	if err := htmlTemplates.ExecuteTemplate(&buf, "notification.html", ctx); err != nil {
		logger.Errorf("Format notification event failed: %s", err)
	}
	m.BodyHtml = buf.String()
	buf.Reset()

	if err := textTemplates.ExecuteTemplate(&buf, "notification.plain", ctx); err != nil {
		logger.Errorf("Format notification event failed: %s", err)
	}
	m.BodyPlain = buf.String()

	return m
}

func (f Formatter) formatMalfunction(sender string, e *events.Malfunction) *Message {
	m := &Message{
		Severity: Fault,
		Subject:  fmt.Sprintf("RMS. Сбой в системе %s: %s", e.System, e.Code),
	}

	ctx := uiContext{
		Title:      m.Subject,
		Time:       time.Unix(e.Timestamp, 0).Local().Format(time.RFC3339),
		Text:       e.Error,
		Sender:     sender,
		Code:       e.Code.String(),
		System:     e.System.String(),
		StackTrace: e.StackTrace,
	}

	var buf bytes.Buffer
	if err := htmlTemplates.ExecuteTemplate(&buf, "malfunction.html", ctx); err != nil {
		logger.Warnf("Format malfunction event failed: %s", err)
	}
	m.BodyHtml = buf.String()
	buf.Reset()

	if err := textTemplates.ExecuteTemplate(&buf, "malfunction.plain", ctx); err != nil {
		logger.Warnf("Format malfunction event failed: %s", err)
	}
	m.BodyPlain = buf.String()

	return m
}

func (f Formatter) formatAlert(sender string, e *events.Alert) *Message {
	m := &Message{}

	m.Severity = Warning
	m.Subject = ""
	switch e.Kind {
	case events.Alert_MotionDetected:
		m.Subject = "Замечено подозрительное движение"
	case events.Alert_CrossLineDetected:
		m.Subject = "Пересечении линии периметра"
	case events.Alert_IntrusionDetected:
		m.Subject = "Нарушение периметра"
	case events.Alert_TamperDetected:
		m.Subject = "Вероятная попытка засветки камеры"
	case events.Alert_GuestDetected:
		m.Severity = Info
		m.Subject = "Гость на входе"
	default:
		m.Subject = "Тревога системы безопасности"
	}

	m.Subject = fmt.Sprintf("RMS. %s, камера %s", m.Subject, e.Camera)

	ts := time.Unix(e.Timestamp, 0).Local()

	ctx := uiContext{
		Title:   m.Subject,
		Time:    ts.Format(time.RFC3339),
		Sender:  prettySender(sender),
		Kind:    e.Kind.String(),
		Channel: e.Camera,
	}

	var buf bytes.Buffer
	if err := htmlTemplates.ExecuteTemplate(&buf, "alert.html", ctx); err != nil {
		logger.Warnf("Format alert event failed: %s", err)
	}
	m.BodyHtml = buf.String()

	buf.Reset()
	if err := textTemplates.ExecuteTemplate(&buf, "alert.plain", ctx); err != nil {
		logger.Warnf("Format alert event failed: %s", err)
	}
	m.BodyPlain = buf.String()

	m.Attachment = &communication.Attachment{
		Type:     communication.Attachment_Photo,
		MimeType: e.ImageMimeType,
		Content:  e.Image,
	}

	// TODO: а как в уведомлении задать кнопку выгрузки архива для tg ?
	return m
}

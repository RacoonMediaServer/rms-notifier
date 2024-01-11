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
	case "rms-notes":
		return "Сервис заметок"
	case "rms-cctv":
		return "Система видеонаблюдения"
	case "rms-bot-client":
		return "Telegram-бот"
	case "rms-library":
		return "Библиотека мультимедия"
	case "rms-web":
		return "Веб-интерфейс"
	case "rms-transcoder":
		return "Транскодер"
	case "rms-backup":
		return "Сервис резервного копирования"
	default:
		return sender
	}
}

func (f Formatter) Format(event interface{}) (*Message, error) {
	switch e := event.(type) {
	case *events.Notification:
		return f.formatNotification(e), nil
	case *events.Malfunction:
		return f.formatMalfunction(e), nil
	case *events.Alert:
		return f.formatAlert(e), nil
	default:
		return nil, fmt.Errorf("unknown event type: %T", e)
	}
}

func (f Formatter) formatNotification(e *events.Notification) *Message {
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
	case events.Notification_BackupComplete:
		m.Severity = Info
		m.Subject = "Резервное копирование завершено"
	case events.Notification_TorrentRemoved:
		// системное уведомление, так что не отправляем
		return nil
	default:
		logger.Errorf("Unknown notification code: %s", e.Kind)
		m.Subject = "Уведомление"
	}

	ctx := uiContext{
		Title:  m.Subject,
		Sender: prettySender(e.Sender),
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

func (f Formatter) formatMalfunction(e *events.Malfunction) *Message {
	m := &Message{
		Severity: Fault,
		Subject:  fmt.Sprintf("RMS. Сбой в системе %s: %s", e.System, e.Code),
	}

	ctx := uiContext{
		Title:      m.Subject,
		Time:       time.Unix(e.Timestamp, 0).Local().Format(time.RFC3339),
		Text:       e.Error,
		Sender:     e.Sender,
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

func (f Formatter) formatAlert(e *events.Alert) *Message {
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
		Sender:  prettySender(e.Sender),
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

	if len(e.Image) != 0 {
		m.Attachment = &communication.Attachment{
			Type:     communication.Attachment_Photo,
			MimeType: e.ImageMimeType,
			Content:  e.Image,
		}
	}

	// TODO: а как в уведомлении задать кнопку выгрузки архива для tg ?
	return m
}

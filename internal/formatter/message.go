package formatter

import "github.com/RacoonMediaServer/rms-packages/pkg/communication"

type Severity int

const (
	Debug Severity = iota
	Info
	Warning
	Fault
)

type Message struct {
	Severity   Severity
	Subject    string
	BodyHtml   string
	BodyPlain  string
	Attachment *communication.Attachment
	Buttons    []*communication.Button
}

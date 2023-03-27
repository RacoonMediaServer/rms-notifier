package service

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-notifier/internal/notifier"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	f         servicemgr.ServiceFactory
	db        Database
	formatter Formatter
	n         Notifier

	enabled atomic.Bool

	mu       sync.Mutex
	settings *rms_notifier.NotifierSettings
}

func (s *Service) Initialize(server server.Server) error {
	settings, err := s.db.LoadSettings(context.Background())
	if err != nil {
		return fmt.Errorf("load settings failed: %w", err)
	}

	s.applySettings(settings)

	if err = s.subscribe(server); err != nil {
		return fmt.Errorf("subscribe to events failed: %w", err)
	}

	return nil
}

func (s *Service) GetSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_notifier.NotifierSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	settings.Enabled = s.settings.Enabled
	settings.Rules = s.settings.Rules
	settings.FilterInterval = s.settings.FilterInterval
	settings.RotationInterval = s.settings.RotationInterval
	return nil
}

func (s *Service) SetSettings(ctx context.Context, settings *rms_notifier.NotifierSettings, empty *emptypb.Empty) error {
	if err := s.db.SaveSettings(ctx, settings); err != nil {
		logger.Errorf("Save settings failed: %s", err)
		return err
	}
	s.applySettings(settings)
	return nil
}

func (s *Service) GetJournal(ctx context.Context, request *rms_notifier.GetJournalRequest, response *rms_notifier.GetJournalResponse) error {
	var from, to *time.Time
	if request.From != nil {
		from = new(time.Time)
		*from = time.Unix(*request.From, 0)
	}
	if request.To != nil {
		to = new(time.Time)
		*to = time.Unix(*request.To, 0)
	}
	evs, err := s.db.LoadEvents(ctx, from, to, uint(request.Limit))
	if err != nil {
		return err
	}

	response.Events = make([]*rms_notifier.Event, 0, len(evs))
	for _, e := range evs {
		response.Events = append(response.Events, &rms_notifier.Event{
			Notification: e.Notification,
			Malfunction:  e.Malfunction,
			Alert:        e.Alert,
			Sender:       e.Sender,
			Timestamp:    e.Timestamp.Unix(),
		})
	}
	return nil
}

func (s *Service) applySettings(settings *rms_notifier.NotifierSettings) {
	logger.Infof("Settings: %+v", settings)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.settings = settings
	s.enabled.Store(settings.Enabled)
	s.n.SetSettings(notifier.Settings{
		Rules: s.settings.Rules,
	})
}

func New(f servicemgr.ServiceFactory, db Database, formatter Formatter, notifier Notifier) *Service {
	return &Service{f: f, db: db, formatter: formatter, n: notifier}
}

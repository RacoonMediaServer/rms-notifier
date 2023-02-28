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
)

type Service struct {
	f         servicemgr.ServiceFactory
	db        Database
	formatter Formatter
	n         Notifier

	enabled atomic.Bool

	mu       sync.Mutex
	settings *rms_notifier.Settings
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

func (s *Service) GetSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_notifier.Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	settings = s.settings
	return nil
}

func (s *Service) SetSettings(ctx context.Context, settings *rms_notifier.Settings, empty *emptypb.Empty) error {
	if err := s.db.SaveSettings(ctx, settings); err != nil {
		logger.Errorf("Save settings failed: %s", err)
		return err
	}
	s.applySettings(settings)
	return nil
}

func (s *Service) GetJournalEvents(ctx context.Context, request *rms_notifier.GetEventsRequest, response *rms_notifier.GetEventsResponse) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) applySettings(settings *rms_notifier.Settings) {
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

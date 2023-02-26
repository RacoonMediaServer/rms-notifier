package service

import (
	"context"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	f servicemgr.ServiceFactory
}

func (s *Service) GetSettings(ctx context.Context, empty *emptypb.Empty, settings *rms_notifier.Settings) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SetSettings(ctx context.Context, settings *rms_notifier.Settings, empty *emptypb.Empty) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetJournalEvents(ctx context.Context, request *rms_notifier.GetEventsRequest, response *rms_notifier.GetEventsResponse) error {
	//TODO implement me
	panic("implement me")
}

func New(f servicemgr.ServiceFactory) *Service {
	return &Service{f: f}
}

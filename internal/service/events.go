package service

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-packages/pkg/pubsub"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

func (s *Service) subscribe(server server.Server) error {
	if err := micro.RegisterSubscriber(pubsub.NotificationTopic, server, s.handleNotification); err != nil {
		return err
	}

	if err := micro.RegisterSubscriber(pubsub.MalfunctionTopic, server, s.handleMalfunction); err != nil {
		return err
	}

	if err := micro.RegisterSubscriber(pubsub.AlertTopic, server, s.handleAlert); err != nil {
		return err
	}

	return nil
}

func (s *Service) handleNotification(ctx context.Context, event events.Notification) error {
	if !s.enabled.Load() {
		return nil
	}

	md, _ := metadata.FromContext(ctx)
	logger.Debugf("Received notification %+v with metadata %+v\n", event, md)

	if err := s.db.StoreEvent(ctx, &event); err != nil {
		logger.Warnf("Store notification event failed: %s", err)
	}

	msg, err := s.formatter.Format(&event)
	if err != nil {
		logger.Errorf("Format event %+v failed: %s", &event, err)
		return nil
	}
	logger.Info(msg.BodyPlain)

	if msg != nil {
		s.n.Notify(pubsub.NotificationTopic, msg)
	}

	return nil
}

func (s *Service) handleMalfunction(ctx context.Context, event events.Malfunction) error {
	if !s.enabled.Load() {
		return nil
	}

	md, _ := metadata.FromContext(ctx)
	logger.Debugf("Received malfunction %+v with metadata %+v\n", event, md)

	if err := s.db.StoreEvent(ctx, &event); err != nil {
		logger.Warnf("Store malfunction event failed: %s", err)
	}

	msg, err := s.formatter.Format(&event)
	if err != nil {
		logger.Errorf("Format event %+v failed: %s", &event, err)
		return nil
	}
	if msg != nil {
		s.n.Notify(pubsub.MalfunctionTopic, msg)
	}

	return nil
}

func (s *Service) handleAlert(ctx context.Context, event events.Alert) error {
	if !s.enabled.Load() {
		return nil
	}

	md, _ := metadata.FromContext(ctx)
	logger.Debugf("Received alert %s from camera %s with metadata %+v\n", event.Kind, event.Camera, md)

	if err := s.db.StoreEvent(ctx, &event); err != nil {
		logger.Warnf("Store alert event failed: %s", err)
	}

	msg, err := s.formatter.Format(&event)
	if err != nil {
		logger.Errorf("Format event %+v failed: %s", &event, err)
		return nil
	}
	logger.Debugf("Formatted: %s", msg.BodyHtml)

	if msg != nil {
		s.n.Notify(pubsub.AlertTopic, msg)
	}

	return nil
}

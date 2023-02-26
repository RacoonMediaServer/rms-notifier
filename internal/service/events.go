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

func (s *Service) Subscribe(server server.Server) error {
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
	md, _ := metadata.FromContext(ctx)
	logger.Debugf("Received notification %+v with metadata %+v\n", event, md)

	return nil
}

func (s *Service) handleMalfunction(ctx context.Context, event events.Malfunction) error {
	md, _ := metadata.FromContext(ctx)
	logger.Debugf("Received malfunction %+v with metadata %+v\n", event, md)

	return nil
}

func (s *Service) handleAlert(ctx context.Context, event events.Alert) error {
	md, _ := metadata.FromContext(ctx)
	logger.Debugf("Received alert %s from camera %s with metadata %+v\n", event.Kind, event.Camera, md)

	return nil
}

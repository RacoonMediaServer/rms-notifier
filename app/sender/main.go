package main

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-packages/pkg/pubsub"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"io/ioutil"
	"runtime"
	"time"
)

func main() {
	var topic string
	var text string
	var kind int
	var system int
	var code int
	var channel string
	var imagePath string

	// create a service
	service := micro.NewService(
		micro.Name("rms-debug.sender"),
		micro.Flags(
			&cli.StringFlag{
				Name:        "topic",
				Usage:       "event topic",
				Required:    false,
				Value:       "rms.notifications",
				DefaultText: "rms.notifications",
				Destination: &topic,
			},
			&cli.StringFlag{
				Name:        "text",
				Usage:       "event text",
				Required:    true,
				Destination: &text,
			},
			&cli.IntFlag{
				Name:        "kind",
				Usage:       "event kind",
				Value:       0,
				Destination: &kind,
			},
			&cli.IntFlag{
				Name:        "system",
				Usage:       "malfunction system",
				Value:       0,
				Destination: &system,
			},
			&cli.IntFlag{
				Name:        "code",
				Usage:       "malfunction code",
				Value:       0,
				Destination: &code,
			},
			&cli.StringFlag{
				Name:        "channel",
				Usage:       "alert channel",
				Destination: &channel,
			},
			&cli.StringFlag{
				Name:        "image",
				Usage:       "attach an image",
				Destination: &imagePath,
			},
		),
	)
	// parse command line
	service.Init()

	logger.Init(logger.WithLevel(logger.TraceLevel))

	image := make([]byte, 0)
	if imagePath != "" {
		var err error
		image, err = ioutil.ReadFile(imagePath)
		if err != nil {
			panic(err)
		}
	}

	pub := pubsub.NewPublisher(service)
	var pkg interface{}

	switch topic {
	case pubsub.NotificationTopic:
		not := &events.Notification{
			ItemTitle: &text,
			Kind:      events.Notification_Kind(kind),
		}

		pkg = not

	case pubsub.MalfunctionTopic:
		b := make([]byte, 8192)
		n := runtime.Stack(b, false)
		s := string(b[:n])

		malf := &events.Malfunction{
			Timestamp:  time.Now().Unix(),
			Error:      text,
			System:     events.Malfunction_System(system),
			Code:       events.Malfunction_Code(code),
			StackTrace: s,
		}

		pkg = malf
	case pubsub.AlertTopic:
		alert := &events.Alert{
			Timestamp: time.Now().Unix(),
			Kind:      events.Alert_Kind(kind),
			Camera:    channel,
			Image:     image,
		}

		pkg = alert

	default:
		panic("unknown topic")
	}

	for {
		if err := pub.Publish(context.Background(), pkg); err != nil {
			logger.Fatalf("publish failed: %s", err)
		}
		<-time.After(time.Second)
	}
}

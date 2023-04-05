package main

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-notifier/internal/config"
	"github.com/RacoonMediaServer/rms-notifier/internal/db"
	"github.com/RacoonMediaServer/rms-notifier/internal/formatter"
	"github.com/RacoonMediaServer/rms-notifier/internal/notifier"
	"github.com/RacoonMediaServer/rms-notifier/internal/sender"
	notifierService "github.com/RacoonMediaServer/rms-notifier/internal/service"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var Version = "v0.0.0"

const serviceName = "rms-notifier"

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	service.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	database, err := db.Connect(config.Config().Database)
	if err != nil {
		logger.Fatalf("Connect to database failed: %s", err)
	}

	f := servicemgr.NewServiceFactory(service)
	n := notifier.New(sender.NewFactory(f, config.Config().Remote, config.Config().Device))
	defer n.Stop()

	srv := notifierService.New(f, database, &formatter.Formatter{}, n)
	if err = srv.Initialize(service.Server()); err != nil {
		logger.Fatalf("Initialize service failed: %s", err)
	}

	// регистрируем хендлеры
	if err = rms_notifier.RegisterRmsNotifierHandler(service.Server(), srv); err != nil {
		logger.Fatalf("Register service failed: %s", err)
	}

	if err = service.Run(); err != nil {
		logger.Fatalf("Run service failed: %s", err)
	}
}

package config

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/pubsub"
	rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"
)

// DefaultSettings is default service settings values
var DefaultSettings = rms_notifier.Settings{
	Enabled: true,
	Rules: map[string]*rms_notifier.Settings_Rules{
		pubsub.NotificationTopic: {
			Rule: []*rms_notifier.Rule{
				{
					Method: rms_notifier.Rule_Telegram,
				},
			},
		},
		pubsub.MalfunctionTopic: {
			Rule: []*rms_notifier.Rule{
				{
					Method: rms_notifier.Rule_Telegram,
				},
			},
		},
		pubsub.AlertTopic: {
			Rule: []*rms_notifier.Rule{
				{
					Method: rms_notifier.Rule_Telegram,
				},
			},
		},
	},
	FilterInterval:   10,
	RotationInterval: 30,
}

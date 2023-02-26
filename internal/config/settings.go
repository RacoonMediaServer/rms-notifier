package config

import rms_notifier "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-notifier"

// DefaultSettings is default service settings values
var DefaultSettings = rms_notifier.Settings{
	Enabled:               true,
	TelegramNotifications: true,
	FilterInterval:        10,
	RotationInterval:      30,
}

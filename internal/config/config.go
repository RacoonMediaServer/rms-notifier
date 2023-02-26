package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Database is credentials to database access
type Database struct {
	// Path to SQLite database
	Path string
}

// Remote is a settings for connection to rms-sender service
type Remote struct {
	Scheme string
	Host   string
	Port   int
	Path   string
}

// Configuration represents entire service configuration
type Configuration struct {
	Database Database
	Remote   Remote
	Device   string
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}

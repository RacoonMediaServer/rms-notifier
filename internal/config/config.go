package config

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-packages/pkg/configuration"
)

// Remote is a settings for connection to rms-sender service
type Remote struct {
	Scheme string
	Host   string
	Port   int
	Path   string
}

// Configuration represents entire service configuration
type Configuration struct {
	Database     string
	Remote       Remote
	APIKey       string
	Mailer       string
	EmailBackend EmailBackend
}

type EmailBackend int

const (
	EmailBackend_RMS EmailBackend = iota
	EmailBackend_TrueNAS
)

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	if err := configuration.Load(configFilePath, &config); err != nil {
		return err
	}

	if config.Mailer != "rms" && config.Mailer != "truenas" {
		return fmt.Errorf("unknown mail server: %s", config.Mailer)
	}

	if config.Mailer == "truenas" {
		config.EmailBackend = EmailBackend_TrueNAS
	}

	return nil
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}

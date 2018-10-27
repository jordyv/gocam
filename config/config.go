package config

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jordyv/gocam/alerting"
	"github.com/spf13/viper"
)

// Config holds all configuration for Gocam
type Config struct {
	Interval            time.Duration
	MaxKeepedImageFiles int
	CameraURL           string
	ImagePath           string
	AlertImagePath      string
	Treshold            int
	AlertHandlers       []string
	LogLevel            log.Level
	HTTPEnabled         bool
	HTTPAddr            string
	MetricsEnabled      bool
	MetricsAddr         string
}

const (
	configInterval            = "interval"
	configImagePath           = "imagePath"
	configAlertImagePath      = "alertImagePath"
	configCameraURL           = "cameraURL"
	configMaxKeepedImageFiles = "maxKeepedImageFiles"
	configVerbose             = "verbose"
	configTreshold            = "treshold"
	configAlertHandlers       = "alertHandlers"
	configAlertHandlerOptions = "alertHandlerOptions"
	configHTTPEnabled         = "http"
	configHTTPAddr            = "httpAddr"
	configMetricsEnabled      = "metrics"
	configMetricsAddr         = "metricsAddr"

	configTelegramOptions = "telegram"
)

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("gocam")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")

	currentDir, _ := os.Getwd()
	viper.SetDefault(configImagePath, fmt.Sprintf("%s/images", currentDir))
	viper.SetDefault(configAlertImagePath, fmt.Sprintf("%s/images/alert", currentDir))
	viper.SetDefault(configInterval, 5*time.Second)
	viper.SetDefault(configMaxKeepedImageFiles, 5)
	viper.SetDefault(configVerbose, false)
	viper.SetDefault(configTreshold, 8)
	viper.SetDefault(configAlertHandlers, []string{})
	viper.SetDefault(configHTTPEnabled, false)
	viper.SetDefault(configHTTPAddr, ":6090")
	viper.SetDefault(configMetricsEnabled, false)
	viper.SetDefault(configMetricsAddr, ":6091")

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

// Build will initialize the configuration from the config file
func (c *Config) Build() {
	c.Interval = viper.GetDuration(configInterval)
	c.MaxKeepedImageFiles = viper.GetInt(configMaxKeepedImageFiles)
	c.CameraURL = viper.GetString(configCameraURL)
	c.ImagePath = viper.GetString(configImagePath)
	c.AlertImagePath = viper.GetString(configAlertImagePath)
	c.Treshold = viper.GetInt(configTreshold)
	c.AlertHandlers = viper.GetStringSlice(configAlertHandlers)
	c.HTTPEnabled = viper.GetBool(configHTTPEnabled)
	c.HTTPAddr = viper.GetString(configHTTPAddr)
	c.MetricsEnabled = viper.GetBool(configMetricsEnabled)
	c.MetricsAddr = viper.GetString(configMetricsAddr)

	if viper.GetBool(configVerbose) {
		c.LogLevel = log.DebugLevel
	} else {
		c.LogLevel = log.InfoLevel
	}
}

// HasAlertHandler checks if a alert handler is enabled
func (c *Config) HasAlertHandler(alertHandlerName string) bool {
	if len(c.AlertHandlers) > 0 {
		for _, h := range c.AlertHandlers {
			if h == alertHandlerName {
				return true
			}
		}
	}
	return false
}

// GetAlertHandlerOptions returns the options for a specific alert handler
func (c *Config) GetAlertHandlerOptions(alertHandlerName string) (interface{}, error) {
	allOptions := viper.GetStringMap(configAlertHandlerOptions)
	if allOptions[alertHandlerName] != nil {
		return allOptions[alertHandlerName], nil
	}
	return nil, fmt.Errorf("couldn't find options for alert handler %s", alertHandlerName)
}

// GetLogAlertHandlerOptions returns the options for the log alert handler
func (c *Config) GetLogAlertHandlerOptions() *alerting.LogAlertHandlerOptions {
	return &alerting.LogAlertHandlerOptions{FilePath: "/dev/stderr"}
}

// GetDiffImageAlertHandlerOptions returns the options for the difference image alert handler
func (c *Config) GetDiffImageAlertHandlerOptions() *alerting.DiffImageAlertHandlerOptions {
	return &alerting.DiffImageAlertHandlerOptions{ImagePath: fmt.Sprintf("%s/%s", c.ImagePath, "diff")}
}

// GetTelegramAlertHandlerOptions returns the options for the telegram alert handler
func (c *Config) GetTelegramAlertHandlerOptions() *alerting.TelegramAlertHandlerOptions {
	viper.SetDefault(configTelegramOptions, map[string]string{"chatID": "", "token": ""})
	telegramConfig := viper.GetStringMapString(configTelegramOptions)
	return &alerting.TelegramAlertHandlerOptions{
		ChatID:   telegramConfig["chatid"],
		BotToken: telegramConfig["token"],
	}
}

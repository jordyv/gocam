package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jordyv/gocam/alerting"
	"github.com/spf13/viper"
)

type Config struct {
	Interval            time.Duration
	MaxKeepedImageFiles int
	CameraURL           string
	ImagePath           string
	Treshold            int
	AlertHandlers       []string
	LogLevel            log.Level
	HTTPEnabled         bool
	HTTPAddr            string
}

const (
	configInterval            = "interval"
	configImagePath           = "imagePath"
	configCameraURL           = "cameraURL"
	configMaxKeepedImageFiles = "maxKeepedImageFiles"
	configVerbose             = "verbose"
	configTreshold            = "treshold"
	configAlertHandlers       = "alertHandlers"
	configAlertHandlerOptions = "alertHandlerOptions"
	configHTTPEnabled         = "http"
	configHTTPAddr            = "httpAddr"

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
	viper.SetDefault(configInterval, 5*time.Second)
	viper.SetDefault(configMaxKeepedImageFiles, 5)
	viper.SetDefault(configVerbose, false)
	viper.SetDefault(configTreshold, 8)
	viper.SetDefault(configAlertHandlers, []string{})
	viper.SetDefault(configHTTPEnabled, false)
	viper.SetDefault(configHTTPAddr, ":6090")

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

func (c *Config) Build() {
	c.Interval = viper.GetDuration(configInterval)
	c.MaxKeepedImageFiles = viper.GetInt(configMaxKeepedImageFiles)
	c.CameraURL = viper.GetString(configCameraURL)
	c.ImagePath = viper.GetString(configImagePath)
	c.Treshold = viper.GetInt(configTreshold)
	c.AlertHandlers = viper.GetStringSlice(configAlertHandlers)
	c.HTTPEnabled = viper.GetBool(configHTTPEnabled)
	c.HTTPAddr = viper.GetString(configHTTPAddr)

	if viper.GetBool(configVerbose) {
		c.LogLevel = log.DebugLevel
	} else {
		c.LogLevel = log.InfoLevel
	}
}

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

func (c *Config) GetAlertHandlerOptions(alertHandlerName string) (interface{}, error) {
	allOptions := viper.GetStringMap(configAlertHandlerOptions)
	if allOptions[alertHandlerName] != nil {
		return allOptions[alertHandlerName], nil
	}
	return nil, errors.New(fmt.Sprintf("couldn't find options for alert handler %s", alertHandlerName))
}

func (c *Config) GetLogAlertHandlerOptions() *alerting.LogAlertHandlerOptions {
	return &alerting.LogAlertHandlerOptions{FilePath: "/dev/stderr"}
}

func (c *Config) GetDiffImageAlertHandlerOptions() *alerting.DiffImageAlertHandlerOptions {
	return &alerting.DiffImageAlertHandlerOptions{ImagePath: fmt.Sprintf("%s/%s", c.ImagePath, "diff")}
}

func (c *Config) GetTelegramAlertHandlerOptions() *alerting.TelegramAlertHandlerOptions {
	viper.SetDefault(configTelegramOptions, map[string]string{"chatID": "", "token": ""})
	telegramConfig := viper.GetStringMapString(configTelegramOptions)
	return &alerting.TelegramAlertHandlerOptions{
		ChatID:   telegramConfig["chatid"],
		BotToken: telegramConfig["token"],
	}
}

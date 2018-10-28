package alerting

import (
	"time"
)

type Alert struct {
	previousImagePath string
	currentImagePath  string
	timestamp         time.Time
	level             int
}

func NewAlert(previousImagePath string, currentImagePath string, level int) *Alert {
	return &Alert{
		previousImagePath: previousImagePath,
		currentImagePath:  currentImagePath,
		level:             level,
		timestamp:         time.Now(),
	}
}

// AlertHandler  Handler for an alert
type AlertHandler interface {
	Notify(alert *Alert)
}

// AlertManager  Handles alertings triggered by the image comparasion
type AlertManager struct {
	handlers []AlertHandler
}

// New  Create a new AlertManager instance
func New(handlers []AlertHandler) *AlertManager {
	return &AlertManager{handlers: handlers}
}

// Notify  Send a notify for an alert to all registered handlers
func (m *AlertManager) Notify(alert *Alert) {
	for _, h := range m.handlers {
		go h.Notify(alert)
	}
}

// AddAlertHandler adds a new alert handler to the manager
func (m *AlertManager) AddAlertHandler(handler AlertHandler) {
	m.handlers = append(m.handlers, handler)
}

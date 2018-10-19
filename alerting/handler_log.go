package alerting

import (
	"fmt"
	"io/ioutil"
	"os"
)

type LogAlertHandlerOptions struct {
	FilePath string
}

// LogAlertHandler  Just print the image path to the stdout when an alert gets triggered
type LogAlertHandler struct {
	Options *LogAlertHandlerOptions
}

func (h *LogAlertHandler) Notify(alert *Alert) {
	logLine := fmt.Sprintf("Alerting for image %s with level %d\n", alert.currentImagePath, alert.level)
	ioutil.WriteFile(h.Options.FilePath, []byte(logLine), os.ModePerm)
}

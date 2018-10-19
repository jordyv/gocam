package alerting

import (
	"fmt"
	"os/exec"
	"time"
)

type DiffImageAlertHandlerOptions struct {
	ImagePath string
	Fuzz      int
}

// DiffImageAlertHandler  Alert handler which will create a difference image from the previous and current alerting image
type DiffImageAlertHandler struct {
	Options *DiffImageAlertHandlerOptions
}

func (h *DiffImageAlertHandler) Notify(alert *Alert) {
	diffImagePath := fmt.Sprintf("%s/%s.jpg", h.Options.ImagePath, alert.timestamp.Format(time.RFC3339))
	fuzz := h.Options.Fuzz
	if fuzz == 0 {
		fuzz = 5
	}
	exec.Command(
		"/usr/bin/compare",
		alert.previousImagePath,
		alert.currentImagePath,
		"-compose", "overlay",
		"-fuzz", fmt.Sprintf("%d%%", fuzz),
		diffImagePath,
	).Run()
}

package alerting

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

var (
	// Max number of recurring alerts in buffer before actually send the message
	recurringBufferCount = 5
)

type TelegramAlertHandlerOptions struct {
	BotToken string
	ChatID   string
}

type TelegramAlertHandler struct {
	Options        *TelegramAlertHandlerOptions
	recurringCount int
}

func (h *TelegramAlertHandler) sendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", h.Options.BotToken, h.Options.ChatID, message)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("Error in Telegram request: %s", string(body))
	}
	return nil
}

func (h *TelegramAlertHandler) Notify(alert *Alert) {
	// Don't send a message for every alert,
	// wait 5 times before sending one.
	if h.recurringCount < recurringBufferCount {
		h.recurringCount++
		return
	}

	h.recurringCount = 0
	message := fmt.Sprintf("Gocam - movement detected!")
	err := h.sendMessage(message)
	if err != nil {
		log.Errorln(err)
	}
}

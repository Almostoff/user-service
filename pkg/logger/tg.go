package logger

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func (a *ApiLogger) InitTG() error {
	var err error
	a.tgBot, err = tb.NewBot(tb.Settings{
		Token:  a.cfg.Logger.TGToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	return err
}

func (a *ApiLogger) SendLogMessage(message string) {
	_, _ = a.tgBot.Send(tb.ChatID(a.cfg.Logger.ChatID), fmt.Sprintf("%s", message))
}

func (a *ApiLogger) makeAlertMessage() (message string) {
	message = "Alert: "
	for i, r := range a.cfg.Logger.AlertUsers {
		message += r
		if i != len(a.cfg.Logger.AlertUsers)-1 {
			message += " "
		}
	}
	return
}

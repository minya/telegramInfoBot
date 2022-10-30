package core

import (
	"fmt"
	"log"
	"time"

	"github.com/minya/telegram"
	"github.com/minya/telegramInfoBot/model"
)

func Run(
	botSettings BotSettings,
	userStorage model.UserStorage,
	process func(upd telegram.Update, h *Handler) interface{},
	notifyCallback UpdateLoopCallback) error {
	botToken := botSettings.GetBotToken()
	ntf := Notifier{
		botToken: botToken,
		storage:  userStorage,
		settings: botSettings.GetNotifierSettings(),
		callback: notifyCallback,
	}
	botApi := telegram.NewApi(botToken)
	ntf.Start(&botApi)
	h := CreateHandler(userStorage, HandlerConfig{process: process})
	listenErr := telegram.StartListen(botToken, 8080, h.handle)
	if nil != listenErr {
		log.Printf("Unable to start listen: %v\n", listenErr)
	}
	return fmt.Errorf("")
}

type NotifierSettings struct {
	UpdateCheckPeriod time.Duration
}

type BotSettings interface {
	IsValid() bool
	GetBotToken() string
	GetNotifierSettings() NotifierSettings
}

package core

import (
	"log"

	"github.com/minya/telegram"
	"github.com/minya/telegramInfoBot/model"
)

type Handler struct {
	Storage model.UserStorage
	Config  HandlerConfig
}

type HandlerConfig struct {
	process func(upd telegram.Update, h *Handler) interface{}
	data    interface{}
}

func CreateHandler(storage model.UserStorage, config HandlerConfig) Handler {
	return Handler{Storage: storage, Config: config}
}

//handle every incoming update
func (h *Handler) handle(upd telegram.Update) interface{} {
	log.Printf("Update: %v\n", upd)
	userID := upd.CallbackQuery.From.Id
	if userID == 0 {
		userID = upd.Message.From.Id
	}

	userInfo, userInfoErr := h.Storage.GetUserInfo(userID)

	if nil != userInfoErr {
		log.Printf("Login not found for user %v. Creating stub.\n", userID)
		h.Storage.SaveUser(userID, &userInfo)
	} else {
		log.Printf("Login for user %v found: %v\n", userID, userInfo.Login)
	}

	return h.Config.process(upd, h)
}

func GetUserID(upd *telegram.Update) int {
	if upd.CallbackQuery.From.Id != 0 {
		return upd.CallbackQuery.From.Id
	}
	return upd.Message.From.Id
}

func GetReplyToChatID(upd *telegram.Update) int {
	chatToReply := upd.Message.Chat.Id
	if chatToReply == 0 {
		chatToReply = upd.CallbackQuery.Message.Chat.Id
	}
	return chatToReply
}

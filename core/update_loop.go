package core

import (
	"log"
	"time"

	"github.com/minya/telegram"
	"github.com/minya/telegramInfoBot/model"
)

type UpdateLoopCallback func(
	id int,
	userInfo *model.UserInfo,
	accountNum string,
	sub *model.SubscriptionInfo,
	api *telegram.Api,
	userStorage model.UserStorage) error

type Notifier struct {
	settings NotifierSettings
	storage  model.UserStorage
	botToken string
	callback UpdateLoopCallback
}

func (n Notifier) Start(api *telegram.Api) {
	go n.updateLoop(api)
}

func (n *Notifier) updateLoop(api *telegram.Api) {
	for true {
		log.Printf("Update...\n")
		subsMap, err := n.storage.GetUsers()
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			for id, userInfo := range subsMap {
				log.Printf("[Update] Check user %v\n", id)
				n.callback(id, &userInfo, "", nil, api, n.storage)
			}
		}
		time.Sleep(n.settings.UpdateCheckPeriod)
	}
}

package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegram-api/internal/infrastructure_layer/router"
)

// 210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk

type TelegramBot struct {
	BotAPI *tgbotapi.BotAPI
	router router.Router
}

func NewTelegramBot(secret string, router router.Router) TelegramBot {
	bot, err := tgbotapi.NewBotAPI(secret)
	if err != nil {
		log.Panic(err)
	}
	return TelegramBot{
		BotAPI: bot,
		router: router,
	}
}

func (t *TelegramBot) StartTelegramServer(debugFlag bool, timeout int) {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)

	t.BotAPI.Debug = debugFlag

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := t.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg, err := t.router.EntryPoint(update)
		if err != nil {
			log.Printf("error when try to send %d", err)
			return
		}

		if _, err = t.BotAPI.Send(msg); err != nil {
			log.Printf("error when try to send %d", err)
			return
		}

	}

}

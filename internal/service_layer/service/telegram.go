package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// 210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk

type TelegramBot struct {
	BotAPI *tgbotapi.BotAPI
}

func NewTelegramBot(secret string) TelegramBot {
	bot, err := tgbotapi.NewBotAPI(secret)
	if err != nil {
		log.Panic(err)
	}
	return TelegramBot{
		BotAPI: bot,
	}
}

func StartTelegramServer(bot *tgbotapi.BotAPI, debugFlag bool, timeout int) {
	log.Printf("Authorized on account %s", bot.Self.UserName)
	
	bot.Debug = debugFlag

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("error when try to send %d", err)
				return
			}
		}
	}

}

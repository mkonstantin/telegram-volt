package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramBot struct {
	BotAPI tgbotapi.BotAPI
}

func NewTelegramBot(botAPI tgbotapi.BotAPI) TelegramBot {
	return TelegramBot{
		BotAPI: botAPI,
	}
}

func main() {
	fmt.Println("sdfs")

	//officeRepo := service.CreateTGBot("210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk")

	//officeRepo, _, _ := InitializeApplication()
	//fmt.Println(officeRepo)

	bot, err := tgbotapi.NewBotAPI("210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

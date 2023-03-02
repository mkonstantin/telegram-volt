package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/scheduler"
)

type TelegramBot struct {
	BotAPI       *tgbotapi.BotAPI
	userSettler  middleware.UserMW
	jobScheduler scheduler.JobsScheduler
	logger       *zap.Logger
}

func NewTelegramBot(botAPI *tgbotapi.BotAPI, router middleware.UserMW, jobScheduler scheduler.JobsScheduler, logger *zap.Logger) TelegramBot {
	return TelegramBot{
		BotAPI:       botAPI,
		userSettler:  router,
		jobScheduler: jobScheduler,
		logger:       logger,
	}
}

func (t *TelegramBot) StartAsyncScheduler() {
	t.jobScheduler.Start()
}

func (t *TelegramBot) StartTelegramServer(debugFlag bool, timeout int) {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)

	t.BotAPI.Debug = debugFlag

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := t.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		msg, err := t.userSettler.EntryPoint(update)
		if err != nil {
			log.Printf("Error handle EntryPoint %d", err)
			continue
		}
		if msg == nil {
			log.Printf("Error handle EntryPoint msg = nil %d", err)
			continue
		}
		if _, err = t.BotAPI.Send(msg); err != nil {
			log.Printf("Error when try to send message %d", err)
			continue
		}
	}
}

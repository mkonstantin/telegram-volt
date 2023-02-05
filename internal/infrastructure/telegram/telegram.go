package telegram

import (
	"fmt"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"telegram-api/internal/app/scheduler"
	"telegram-api/internal/infrastructure/router"
	"time"
)

type TelegramBot struct {
	BotAPI *tgbotapi.BotAPI
	router router.Router
	logger *zap.Logger
	worker scheduler.Worker
}

func NewTelegramBot(secret string, router router.Router, worker scheduler.Worker, logger *zap.Logger) TelegramBot {
	bot, err := tgbotapi.NewBotAPI(secret)
	if err != nil {
		log.Panic(err)
	}
	return TelegramBot{
		BotAPI: bot,
		router: router,
		worker: worker,
		logger: logger,
	}
}

func (t *TelegramBot) StartAsyncScheduler() {
	s := gocron.NewScheduler(time.FixedZone("UTC+6", 6*60*60))
	_, err := s.Every(1).
		Week().
		At("14:30").
		Weekday(time.Monday).
		Weekday(time.Tuesday).
		Weekday(time.Wednesday).
		Weekday(time.Thursday).
		Weekday(time.Friday).
		Do(func() {
			fmt.Println("do work 1")
			err := t.worker.CleanTables()
			if err != nil {
				t.logger.Error("gocron execution error", zap.Error(err))
			}
		})
	if err != nil {
		t.logger.Error("gocron create error", zap.Error(err))
	}
	s.StartAsync()
}

func (t *TelegramBot) StartTelegramServer(debugFlag bool, timeout int) {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)

	t.BotAPI.Debug = debugFlag

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := t.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		msg, err := t.router.MainEntryPoint(update)
		if err != nil {
			log.Printf("Error handle MainEntryPoint %d", err)
			continue
		}
		if msg == nil {
			log.Printf("Error handle MainEntryPoint msg = nil %d", err)
			continue
		}
		if _, err = t.BotAPI.Send(msg); err != nil {
			log.Printf("Error when try to send message %d", err)
			continue
		}
	}
}

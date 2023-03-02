//go:build wireinject
// +build wireinject

package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/telegram"
)

func InitializeApplication(secret string, cfg config.AppConfig, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		provideTelegramAPI,
		dbSet,
		repositorySet,
		middleware.NewUserMW,
		telegram.NewTelegramBot,
		menuSet,
		handlerSet,
		servicesSet,
		jobsSet,
		formSet,
	)
	return telegram.TelegramBot{}, nil, nil
}

func provideTelegramAPI(secret string, logger *zap.Logger) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(secret)
	if err != nil {
		logger.Panic("err telegram api init", zap.Error(err))
	}
	return bot
}

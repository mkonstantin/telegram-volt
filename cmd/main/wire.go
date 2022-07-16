//go:build wireinject
// +build wireinject

package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/wire"
)

func InitializeApplication(botAPI tgbotapi.BotAPI) (TelegramBot, func(), error) {
	wire.Build(
		NewTelegramBot,
		//repositorySet,
	)
	return TelegramBot{}, nil, nil
}

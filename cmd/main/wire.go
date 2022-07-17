//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"telegram-api/internal/infrastructure_layer/telegram"
)

func InitializeApplication(secret string) (telegram.TelegramBot, func(), error) {
	wire.Build(
		telegram.NewTelegramBot,
		//repositorySet,
	)
	return telegram.TelegramBot{}, nil, nil
}

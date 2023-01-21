//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure_layer/hundlers"
	"telegram-api/internal/infrastructure_layer/router"
	"telegram-api/internal/infrastructure_layer/telegram"
)

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		router.NewRouter,
		telegram.NewTelegramBot,
		hundlers.NewCommandHandler,
		hundlers.NewCustomMessageHandler,
		hundlers.NewInlineMessageHandler,
		//hundlers.NewPrimaryHundler,
		//repositorySet,
		//dbSet,
	)
	return telegram.TelegramBot{}, nil, nil
}

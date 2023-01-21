//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/handler"
	"telegram-api/internal/infrastructure/router"
	"telegram-api/internal/infrastructure/telegram"
)

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		router.NewRouter,
		telegram.NewTelegramBot,
		handler.NewCommandHandler,
		handler.NewCustomMessageHandler,
		handler.NewInlineMessageHandler,
		//handler.NewPrimaryHundler,
		//repositorySet,
		//dbSet,
	)
	return telegram.TelegramBot{}, nil, nil
}

//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/router"
	"telegram-api/internal/infrastructure/telegram"
)

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		dbSet,
		repositorySet,
		router.NewRouter,
		telegram.NewTelegramBot,
		handlerSet,
		servicesSet,
		jobsSet,
	)
	return telegram.TelegramBot{}, nil, nil
}

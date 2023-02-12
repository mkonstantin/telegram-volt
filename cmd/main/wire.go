//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/telegram"
)

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		dbSet,
		repositorySet,
		middleware.NewUserMW,
		telegram.NewTelegramBot,
		handlerSet,
		servicesSet,
		jobsSet,
	)
	return telegram.TelegramBot{}, nil, nil
}

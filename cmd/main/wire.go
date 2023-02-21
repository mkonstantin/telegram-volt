//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/telegram"
)

func InitializeApplication(secret string, cfg config.AppConfig, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
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

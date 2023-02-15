//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/former"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/telegram"
)

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		dbSet,
		repositorySet,
		former.NewMessageFormer,
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

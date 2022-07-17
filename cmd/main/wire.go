//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure_layer/telegram"
)

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	wire.Build(
		telegram.NewTelegramBot,
		dbSet,
		repositorySet,
	)
	return telegram.TelegramBot{}, nil, nil
}

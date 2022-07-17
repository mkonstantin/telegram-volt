//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"telegram-api/internal/service_layer/service"
)

func InitializeApplication(secret string) (service.TelegramBot, func(), error) {
	wire.Build(
		service.NewTelegramBot,
		//repositorySet,
	)
	return service.TelegramBot{}, nil, nil
}

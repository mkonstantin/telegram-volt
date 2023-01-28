// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler"
	"telegram-api/internal/infrastructure/repo"
	"telegram-api/internal/infrastructure/router"
	"telegram-api/internal/infrastructure/telegram"
)

// Injectors from wire.go:

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	customMessageHandler := handler.NewCustomMessageHandler(logger)
	messageFormer := handler.NewMessageFormer(logger)
	contextContext := context.Background()
	connection, cleanup := provideDBConnection(contextContext, logger)
	userRepository := repo.NewUserRepository(connection)
	officeRepository := repo.NewOfficeRepository(connection)
	seatRepository := repo.NewSeatRepository(connection)
	userService := usecase.NewUserService(userRepository, officeRepository, seatRepository, logger)
	commandHandler := handler.NewCommandHandler(messageFormer, userService, logger)
	inlineMessageHandler := handler.NewInlineMessageHandler(messageFormer, userService, logger)
	routerRouter := router.NewRouter(customMessageHandler, commandHandler, inlineMessageHandler, logger)
	telegramBot := telegram.NewTelegramBot(secret, routerRouter)
	return telegramBot, func() {
		cleanup()
	}, nil
}

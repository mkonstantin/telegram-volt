// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"go.uber.org/zap"
	"telegram-api/internal/app/scheduler"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/handler"
	"telegram-api/internal/infrastructure/repo"
	"telegram-api/internal/infrastructure/router"
	"telegram-api/internal/infrastructure/telegram"
)

// Injectors from wire.go:

func InitializeApplication(secret string, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	contextContext := context.Background()
	connection, cleanup := provideDBConnection(contextContext, logger)
	userRepository := repo.NewUserRepository(connection)
	customMessageHandler := handler.NewCustomMessageHandler(logger)
	messageFormer := handler.NewMessageFormer(logger)
	officeRepository := repo.NewOfficeRepository(connection)
	bookSeatRepository := repo.NewBookSeatRepository(connection)
	userService := usecase.NewUserService(userRepository, officeRepository, bookSeatRepository, logger)
	commandHandler := handler.NewCommandHandler(messageFormer, userService, logger)
	inlineMessageHandler := handler.NewInlineMessageHandler(messageFormer, userService, logger)
	routerRouter := router.NewRouter(userRepository, customMessageHandler, commandHandler, inlineMessageHandler, logger)
	worker := scheduler.NewWorker(logger)
	telegramBot := telegram.NewTelegramBot(secret, routerRouter, worker, logger)
	return telegramBot, func() {
		cleanup()
	}, nil
}

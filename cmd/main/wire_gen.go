// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/handler"
	"telegram-api/internal/app/informer"
	"telegram-api/internal/app/menu"
	"telegram-api/internal/app/scheduler"
	"telegram-api/internal/app/scheduler/job"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/repo"
	"telegram-api/internal/infrastructure/router"
	"telegram-api/internal/infrastructure/telegram"
)

// Injectors from wire.go:

func InitializeApplication(secret string, cfg config.AppConfig, logger *zap.Logger) (telegram.TelegramBot, func(), error) {
	botAPI := provideTelegramAPI(secret, logger)
	contextContext := context.Background()
	connection, cleanup := provideDBConnection(contextContext, cfg, logger)
	userRepository := repo.NewUserRepository(connection)
	workDateRepository := repo.NewWorkDateRepository(connection)
	bookSeatRepository := repo.NewBookSeatRepository(connection)
	officeRepository := repo.NewOfficeRepository(connection)
	officeMenuForm := form.NewOfficeMenuForm(logger)
	officeMenu := menu.NewOfficeMenu(workDateRepository, bookSeatRepository, officeRepository, officeMenuForm, logger)
	officeListForm := form.NewOfficeListForm(logger)
	officeListMenu := menu.NewOfficeListMenu(officeRepository, officeListForm, logger)
	start := handler.NewStartHandle(officeMenu, officeListMenu, logger)
	userService := usecase.NewUserService(userRepository, officeRepository, bookSeatRepository, logger)
	officeList := handler.NewOfficeListHandle(userService, officeMenu, logger)
	infoMenuForm := form.NewInfoMenuForm(logger)
	informerService := informer.NewInformer(botAPI, infoMenuForm, userRepository, bookSeatRepository, logger)
	dateMenuForm := form.NewDateMenutForm(logger)
	dateMenu := menu.NewDateMenu(workDateRepository, officeRepository, bookSeatRepository, dateMenuForm, logger)
	handlerOfficeMenu := handler.NewOfficeMenuHandle(informerService, userService, dateMenu, officeListMenu, logger)
	ownSeatForm := form.NewOwnSeatForm(logger)
	ownSeatMenu := menu.NewOwnSeatMenu(ownSeatForm, logger)
	freeSeatForm := form.NewFreeSeatForm(logger)
	freeSeatMenu := menu.NewFreeSeatMenu(bookSeatRepository, freeSeatForm, logger)
	seatList := handler.NewSeatListHandle(bookSeatRepository, dateMenu, ownSeatMenu, freeSeatMenu, logger)
	seatListForm := form.NewSeatListForm(logger)
	seatListMenu := menu.NewSeatListMenu(bookSeatRepository, seatListForm, logger)
	handlerOwnSeatMenu := handler.NewOwnSeatMenuHandle(informerService, userService, bookSeatRepository, seatListMenu, logger)
	handlerFreeSeatMenu := handler.NewFreeSeatMenuHandle(userService, bookSeatRepository, seatListMenu, logger)
	handlerDateMenu := handler.NewDateMenuHandle(seatListMenu, officeMenu, logger)
	infoMenu := handler.NewInfoMenuHandle(seatListMenu, officeMenu, logger)
	routerRouter := router.NewRouter(start, officeList, handlerOfficeMenu, seatList, handlerOwnSeatMenu, handlerFreeSeatMenu, handlerDateMenu, infoMenu, logger)
	userMW := middleware.NewUserMW(userRepository, routerRouter, logger)
	hourlyJob := job.NewHourlyJob(informerService, officeRepository, workDateRepository, logger)
	dateJob := job.NewDateJob(workDateRepository, logger)
	seatRepository := repo.NewSeatRepository(connection)
	seatJob := job.NewSeatsJob(officeRepository, workDateRepository, bookSeatRepository, seatRepository, logger)
	jobsScheduler := scheduler.NewJobsScheduler(hourlyJob, dateJob, seatJob, logger)
	telegramBot := telegram.NewTelegramBot(botAPI, userMW, jobsScheduler, logger)
	return telegramBot, func() {
		cleanup()
	}, nil
}

// wire.go:

func provideTelegramAPI(secret string, logger *zap.Logger) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(secret)
	if err != nil {
		logger.Panic("err telegram api init", zap.Error(err))
	}
	return bot
}

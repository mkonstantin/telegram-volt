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
	officeRepository := repo.NewOfficeRepository(connection)
	bookSeatRepository := repo.NewBookSeatRepository(connection)
	userService := usecase.NewUserService(userRepository, officeRepository, bookSeatRepository, logger)
	workDateRepository := repo.NewWorkDateRepository(connection)
	officeMenuForm := form.NewOfficeMenuForm(logger)
	officeMenu := menu.NewOfficeMenu(userService, workDateRepository, bookSeatRepository, officeRepository, officeMenuForm, logger)
	officeListForm := form.NewOfficeListForm(logger)
	officeListMenu := menu.NewOfficeListMenu(officeRepository, officeListForm, logger)
	start := handler.NewStartHandle(officeMenu, officeListMenu, logger)
	officeList := handler.NewOfficeListHandle(officeMenu, logger)
	infoMenuForm := form.NewInfoMenuForm(logger)
	sender := informer.NewSender(botAPI, infoMenuForm, logger)
	informerService := informer.NewInformer(botAPI, infoMenuForm, userRepository, bookSeatRepository, sender, logger)
	dateMenuForm := form.NewDateMenutForm(logger)
	dateMenu := menu.NewDateMenu(workDateRepository, officeRepository, bookSeatRepository, dateMenuForm, cfg, logger)
	handlerOfficeMenu := handler.NewOfficeMenuHandle(informerService, userService, dateMenu, officeMenu, officeListMenu, logger)
	ownSeatForm := form.NewOwnSeatForm(logger)
	ownSeatMenu := menu.NewOwnSeatMenu(ownSeatForm, logger)
	holdSeatForm := form.NewHoldSeatForm(logger)
	holdSeatMenu := menu.NewHoldSeatMenu(holdSeatForm, logger)
	freeSeatForm := form.NewFreeSeatForm(logger)
	freeSeatMenu := menu.NewFreeSeatMenu(bookSeatRepository, freeSeatForm, cfg, logger)
	seatList := handler.NewSeatListHandle(bookSeatRepository, dateMenu, ownSeatMenu, holdSeatMenu, freeSeatMenu, cfg, logger)
	seatListForm := form.NewSeatListForm(logger)
	seatListMenu := menu.NewSeatListMenu(userService, bookSeatRepository, seatListForm, logger)
	handlerOwnSeatMenu := handler.NewOwnSeatMenuHandle(officeMenu, informerService, userService, bookSeatRepository, seatListMenu, logger)
	handlerFreeSeatMenu := handler.NewFreeSeatMenuHandle(officeMenu, userService, bookSeatRepository, seatListMenu, logger)
	handlerDateMenu := handler.NewDateMenuHandle(seatListMenu, officeMenu, officeRepository, logger)
	infoMenu := handler.NewInfoMenuHandle(bookSeatRepository, seatListMenu, officeMenu, logger)
	handlerHoldSeatMenu := handler.NewHoldSeatMenuHandle(userService, bookSeatRepository, seatListMenu, logger)
	routerRouter := router.NewRouter(cfg, start, officeList, handlerOfficeMenu, seatList, handlerOwnSeatMenu, handlerFreeSeatMenu, handlerDateMenu, infoMenu, handlerHoldSeatMenu, logger)
	userMW := middleware.NewUserMW(userRepository, routerRouter, logger)
	hourlyJob := job.NewHourlyJob(informerService, officeRepository, workDateRepository, bookSeatRepository, logger)
	dateJob := job.NewDateJob(workDateRepository, logger)
	seatRepository := repo.NewSeatRepository(connection)
	seatJob := job.NewSeatsJob(officeRepository, workDateRepository, bookSeatRepository, seatRepository, logger)
	jobsScheduler := scheduler.NewJobsScheduler(hourlyJob, dateJob, seatJob, logger)
	telegramBot := telegram.NewTelegramBot(botAPI, userMW, jobsScheduler, sender, logger)
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

package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/service"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/router"
	"telegram-api/internal/infrastructure/scheduler"
)

var servicesSet = wire.NewSet(
	router.NewRouter,
	usecase.NewUserService,
	scheduler.NewJobsScheduler,
	service.NewInformer,
)

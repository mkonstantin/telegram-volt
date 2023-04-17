package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/informer"
	"telegram-api/internal/app/scheduler"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/router"
)

var servicesSet = wire.NewSet(
	router.NewRouter,
	usecase.NewUserService,
	scheduler.NewJobsScheduler,
	informer.NewInformer,
	informer.NewSender,
)

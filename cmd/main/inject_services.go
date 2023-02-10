package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/scheduler"
	"telegram-api/internal/infrastructure/service"
)

var servicesSet = wire.NewSet(
	usecase.NewUserService,
	scheduler.NewJobsScheduler,
	service.NewTimeHelper,
)

package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/scheduler"
)

var servicesSet = wire.NewSet(
	usecase.NewUserService,
	scheduler.NewJobsScheduler,
)

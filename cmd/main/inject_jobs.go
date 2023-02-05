package main

import (
	"github.com/google/wire"
	"telegram-api/internal/infrastructure/scheduler/handler"
)

var jobsSet = wire.NewSet(
	handler.NewOfficeJob,
)

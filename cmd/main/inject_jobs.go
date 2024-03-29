package main

import (
	"github.com/google/wire"
	"telegram-api/internal/app/scheduler/job"
)

var jobsSet = wire.NewSet(
	job.NewHourlyJob,
	job.NewDateJob,
	job.NewSeatsJob,
)

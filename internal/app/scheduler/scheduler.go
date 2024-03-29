package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"telegram-api/internal/app/scheduler/job"
)

type jobsSchedulerImpl struct {
	hourlyJob job.HourlyJob
	dateJob   job.DateJob
	seatJob   job.SeatJob
	logger    *zap.Logger
}

type JobsScheduler interface {
	StartFillWorkDates()
	StartHourlyJob()
}

func NewJobsScheduler(hourlyJob job.HourlyJob,
	dateJob job.DateJob,
	seatJob job.SeatJob,
	logger *zap.Logger) JobsScheduler {

	return &jobsSchedulerImpl{
		hourlyJob: hourlyJob,
		dateJob:   dateJob,
		seatJob:   seatJob,
		logger:    logger,
	}
}

// StartFillWorkDates Start FillWork Cron Scheduler

func (w *jobsSchedulerImpl) StartFillWorkDates() {
	w.logger.Info("Starting Fill Work Dates and Seats scheduled job")

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).
		Day().
		At("03:15").
		Do(func() {
			w.logger.Info("gocron start Date & Seat Jobs")

			//err := w.dateJob.CheckAndSetDates()
			//if err != nil {
			//	w.logger.Error("gocron execution DateJobs error", zap.Error(err))
			//}
			//
			//err = w.seatJob.SetSeats()
			//if err != nil {
			//	w.logger.Error("gocron execution SeatJob error", zap.Error(err))
			//}
		})

	if err != nil {
		w.logger.Error("gocron create Date & Seat Jobs error", zap.Error(err))
		return
	}

	s.StartImmediately()
	s.StartAsync()

	w.logger.Info("Successfully started scheduled job: DateJobs, SeatJob")
	return
}

func (w *jobsSchedulerImpl) StartHourlyJob() {
	today := time.Now()
	startedAt := time.Date(today.Year(), today.Month(), today.Day(), today.Hour()+1, 0, 0, 0, time.UTC)
	w.logger.Info(fmt.Sprintf("Hourly Job startedAt: %s", startedAt.String()))

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).
		Hour().
		StartAt(startedAt).
		Do(func() {
			w.logger.Info("gocron start Hourly Job")

			err := w.hourlyJob.StartSchedule()
			if err != nil {
				w.logger.Error("gocron execution Hourly Job error", zap.Error(err))
			}
		})
	if err != nil {
		w.logger.Error("gocron create Hourly Job error", zap.Error(err))
	}

	s.StartImmediately()
	s.StartAsync()

	w.logger.Info("Successfully started scheduled job: Hourly Job")
}

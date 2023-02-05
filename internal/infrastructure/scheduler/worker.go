package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/scheduler/handler"
	"time"
)

type jobsSchedulerImpl struct {
	officeRepo interfaces.OfficeRepository
	officeJobs handler.OfficeJob
	logger     *zap.Logger
}

type JobsScheduler interface {
	Start()
}

func NewJobsScheduler(officeJobs handler.OfficeJob, logger *zap.Logger) JobsScheduler {
	return &jobsSchedulerImpl{
		officeJobs: officeJobs,
		logger:     logger,
	}
}

func (w *jobsSchedulerImpl) Start() {
	w.logger.Info("Starting scheduler jobs")

	s := gocron.NewScheduler(time.FixedZone("UTC+6", 6*60*60))
	_, err := s.Every(1).
		Week().
		At("14:30").
		Weekday(time.Monday).
		Weekday(time.Tuesday).
		Weekday(time.Wednesday).
		Weekday(time.Thursday).
		Weekday(time.Friday).
		Do(func() {
			fmt.Println("do work 1")
			err := w.officeJobs.BeginJob()
			if err != nil {
				w.logger.Error("gocron execution error", zap.Error(err))
			}
		})
	if err != nil {
		w.logger.Error("gocron create error", zap.Error(err))
	}
	s.StartAsync()
}

func (w *jobsSchedulerImpl) CleanTables() error {

	return nil
}

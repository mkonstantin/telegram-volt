package scheduler

import (
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"telegram-api/internal/app/scheduler/job"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

type jobsSchedulerImpl struct {
	officeRepo interfaces.OfficeRepository
	dateJob    job.DateJob
	seatJob    job.SeatJob
	logger     *zap.Logger
}

type JobsScheduler interface {
	StartFillWorkDates()
	StartFillSeats()
}

func NewJobsScheduler(officeRepo interfaces.OfficeRepository,
	dateJob job.DateJob,
	seatJob job.SeatJob,
	logger *zap.Logger) JobsScheduler {

	return &jobsSchedulerImpl{
		officeRepo: officeRepo,
		dateJob:    dateJob,
		seatJob:    seatJob,
		logger:     logger,
	}
}

// StartFillWorkDates Start FillWork Cron Scheduler

func (w *jobsSchedulerImpl) StartFillWorkDates() {
	w.logger.Info("Starting Fill Work Dates scheduled job")

	err := w.startDateJob()
	if err != nil {
		w.logger.Error("Job scheduler get error when try starting DateJob", zap.Error(err))
		return
	}
}

func (w *jobsSchedulerImpl) startDateJob() error {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).
		Day().
		At("03:00").
		Do(func() {
			w.logger.Info("gocron start CheckAndSetDates")

			err := w.dateJob.CheckAndSetDates()
			if err != nil {
				w.logger.Error("gocron execution DateJobs error", zap.Error(err))
			}
		})
	if err != nil {
		w.logger.Error("gocron create DateJobs error", zap.Error(err))
		return err
	}

	s.StartImmediately()
	s.StartAsync()

	w.logger.Info("Successfully started scheduled job: DateJobs")
	return nil
}

// StartFillSeats Start FillSeats Cron Scheduler

func (w *jobsSchedulerImpl) StartFillSeats() {
	w.logger.Info("Starting Fill Seats scheduled job")

	err := w.startSeatsJob()
	if err != nil {
		w.logger.Error("Job scheduler get error when try starting SeatsJob", zap.Error(err))
		return
	}
}

func (w *jobsSchedulerImpl) startSeatsJob() error {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).
		Day().
		At("03:30").
		Do(func() {
			w.logger.Info("gocron start SetSeats")

			err := w.seatJob.SetSeats()
			if err != nil {
				w.logger.Error("gocron execution SeatsJob error", zap.Error(err))
			}
		})
	if err != nil {
		w.logger.Error("gocron create SeatsJob error", zap.Error(err))
		return err
	}

	s.StartImmediately()
	s.StartAsync()

	w.logger.Info("Successfully started scheduled job: SeatsJob")
	return nil
}

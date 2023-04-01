package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"telegram-api/internal/app/scheduler/job"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

type jobsSchedulerImpl struct {
	officeRepo interfaces.OfficeRepository
	officeJobs job.OfficeJob
	dateJob    job.DateJob
	logger     *zap.Logger
}

type JobsScheduler interface {
	StartFillWorkDates()
	StartEnableBook()
}

func NewJobsScheduler(officeRepo interfaces.OfficeRepository, officeJobs job.OfficeJob, dateJob job.DateJob, logger *zap.Logger) JobsScheduler {
	return &jobsSchedulerImpl{
		officeRepo: officeRepo,
		officeJobs: officeJobs,
		dateJob:    dateJob,
		logger:     logger,
	}
}

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
		At("12:30").
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

func (w *jobsSchedulerImpl) StartEnableBook() {
	w.logger.Info("Starting Enable Book scheduled job")

	offices, err := w.officeRepo.GetAll()
	if err != nil {
		w.logger.Error("Scheduler jobs get all offices error", zap.Error(err))
		return
	}

	for _, office := range offices {
		err = w.createForOffice(office)
		if err != nil {
			w.logger.Error("Scheduler jobs set office error", zap.Error(err))
			return
		}
	}
}

func (w *jobsSchedulerImpl) createForOffice(office *model.Office) error {
	location, err := time.LoadLocation(office.TimeZone)
	if err != nil {
		return err
	}

	s := gocron.NewScheduler(location)
	_, err = s.Every(3).
		Hours().
		Do(func() {
			w.logger.Info("gocron start CheckAndSetDates")

			year, week := time.Now().ISOWeek()
			err = w.officeJobs.SetSeatsForAllWeek(office.ID, year, week)
			if err != nil {
				w.logger.Error("gocron execution error", zap.Error(err))
			}
		})

	if err != nil {
		w.logger.Error("gocron create error", zap.Error(err))
	}
	s.StartImmediately()
	s.StartAsync()

	w.logger.Info(fmt.Sprintf("Successfully started async scheduled jobs for %s", office.Name))
	return nil
}

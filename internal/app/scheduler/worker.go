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
	logger     *zap.Logger
}

type JobsScheduler interface {
	Start()
}

func NewJobsScheduler(officeRepo interfaces.OfficeRepository, officeJobs job.OfficeJob, logger *zap.Logger) JobsScheduler {
	return &jobsSchedulerImpl{
		officeRepo: officeRepo,
		officeJobs: officeJobs,
		logger:     logger,
	}
}

func (w *jobsSchedulerImpl) Start() {
	w.logger.Info("Starting scheduler jobs")

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
	_, err = s.Every(1).
		Day().
		At("22:00").
		Do(func() {
			year, week := time.Now().ISOWeek()

			err = w.officeJobs.SetSeatsWeek(office.ID, year, week)
			if err != nil {
				w.logger.Error("gocron execution error", zap.Error(err))
			}
		})

	if err != nil {
		w.logger.Error("gocron create error", zap.Error(err))
	}
	s.StartImmediately()
	s.StartAsync()

	w.logger.Info(fmt.Sprintf("Success start async scheduled jobs for %s", office.Name))
	return nil
}

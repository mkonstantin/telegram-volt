package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"telegram-api/internal/domain/model"
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

func NewJobsScheduler(officeRepo interfaces.OfficeRepository, officeJobs handler.OfficeJob, logger *zap.Logger) JobsScheduler {
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
		//Minute().
		Week().
		At("14:30").
		Weekday(time.Sunday).
		Weekday(time.Monday).
		Weekday(time.Tuesday).
		Weekday(time.Wednesday).
		Weekday(time.Thursday).
		Do(func() {
			err = w.officeJobs.SetNewSeatList(office.ID, location)
			if err != nil {
				w.logger.Error("gocron execution error", zap.Error(err))
			}
		})
	if err != nil {
		w.logger.Error("gocron create error", zap.Error(err))
	}
	s.StartAsync()

	w.logger.Info(fmt.Sprintf("Success start async scheduled jobs for %s", office.Name))
	return nil
}

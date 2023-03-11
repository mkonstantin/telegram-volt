package job

import (
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"time"
)

const (
	plusDaysAmount = 10
	addDaysAmount  = 30
)

type dateJobsImpl struct {
	dateRepo interfaces.WorkDateRepository
	logger   *zap.Logger
}

type DateJob interface {
	CheckAndSetDates() error
}

func NewDateJob(dateRepo interfaces.WorkDateRepository, logger *zap.Logger) DateJob {
	return &dateJobsImpl{
		dateRepo: dateRepo,
		logger:   logger,
	}
}

func (o *dateJobsImpl) CheckAndSetDates() error {

	last, err := o.dateRepo.GetLastByDate()
	if err != nil {
		return err
	}

	todayPlus10Days := helper.TodayPlusUTC(plusDaysAmount)

	if last == nil {
		return o.addDays(helper.TodayZeroTimeUTC())
	}

	// добавляем даты за 10 дней до конца текущих
	if todayPlus10Days.After(last.WorkDate) {
		date := last.WorkDate.AddDate(0, 0, 1)
		return o.addDays(date)
	}

	return nil
}

func (o *dateJobsImpl) addDays(startDate time.Time) error {
	for i := 0; i < addDaysAmount; i++ {
		nextDate := startDate.AddDate(0, 0, i)
		err := o.dateRepo.InsertDate(nextDate)
		if err != nil {
			o.logger.Error("error while add dates", zap.Error(err))
			return err
		}
	}
	return nil
}

package job

import (
	"context"
	"go.uber.org/zap"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/middleware"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
	"time"
)

const (
	checkDays       = 20
	totalDaysAmount = 30
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

	ctx := context.Background()
	ctx = middleware.SetNewTrace(ctx, "DateJob", 100200300)
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	last, err := o.dateRepo.GetLastByDate(ctx)
	if err != nil {
		return err
	}

	// Если дат нет, то сразу добавляем totalDaysAmount дней
	if last == nil {
		return o.addDays(ctx, helper.TodayZeroTimeUTC(), totalDaysAmount)
	}

	// Замеряем сколько дней до лимита даты
	today := helper.TodayZeroTimeUTC()
	duration := last.Date.Sub(today)
	days := int(duration.Hours()/24) + 1

	// Если меньше или равно checkDays, то прибавляем разницу чтобы всегда было + 20-30 дней
	if days <= checkDays {
		date := last.Date.AddDate(0, 0, 1)
		return o.addDays(ctx, date, totalDaysAmount-days)
	}

	return nil
}

func (o *dateJobsImpl) addDays(ctx context.Context, startDate time.Time, daysAmount int) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	for i := 0; i < daysAmount; i++ {
		nextDate := startDate.AddDate(0, 0, i)
		err := o.dateRepo.InsertDate(ctx, nextDate)
		if err != nil {
			o.logger.Error("error while add dates", zap.Error(err))
			return err
		}
	}
	return nil
}

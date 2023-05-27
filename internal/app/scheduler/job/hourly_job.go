package job

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/app/informer"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/common"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type hourlyJobImpl struct {
	informerService informer.InformerService
	officeRepo      interfaces.OfficeRepository
	workDateRepo    interfaces.WorkDateRepository
	bookSeatRepo    interfaces.BookSeatRepository
	logger          *zap.Logger
}

type HourlyJob interface {
	StartSchedule() error
}

func NewHourlyJob(
	informerService informer.InformerService,
	officeRepo interfaces.OfficeRepository,
	workDateRepo interfaces.WorkDateRepository,
	bookSeatRepo interfaces.BookSeatRepository,
	logger *zap.Logger) HourlyJob {
	return &hourlyJobImpl{
		informerService: informerService,
		officeRepo:      officeRepo,
		workDateRepo:    workDateRepo,
		bookSeatRepo:    bookSeatRepo,
		logger:          logger,
	}
}

func (h *hourlyJobImpl) StartSchedule() error {
	ctx := context.Background()

	today := helper.TodayZeroTimeUTC()
	todayPlus2 := helper.TodayPlusUTC(2)

	// переводим все прошедшие даты в статус Done
	err := h.workDateRepo.DoneAllPastByDate(today.String())
	if err != nil {
		h.logger.Error("HourlyJob workDateRepo.DoneAllPastByDate error", zap.Error(err))
		return err
	}

	// получаем сегодня и завтра
	dates, err := h.workDateRepo.FindByDates(today.String(), todayPlus2.String())
	if err != nil {
		h.logger.Error("HourlyJob workDateRepo.FindByDatesAndStatus error", zap.Error(err))
		return err
	}

	if len(dates) == 0 {
		return common.ErrDateSetBookSeatsNotFound
	}

	var todayWorkDate model.WorkDate
	var tomorrowWorkDate model.WorkDate

	if len(dates) > 0 {
		todayWorkDate = dates[0]
	}
	if len(dates) > 1 {
		tomorrowWorkDate = dates[1]
	}

	// в любом случае открываем запись на сегодня
	err = h.openTodayAccept(ctx, todayWorkDate)
	if err != nil {
		h.logger.Error("HourlyJob h.openTodayAccept error", zap.Error(err))
		return err
	}

	offices, err := h.officeRepo.GetAll()
	if err != nil {
		h.logger.Error("HourlyJob officeRepo.GetAll error", zap.Error(err))
		return err
	}

	// чекаем что сегодня рабочий день закончился и закрываем его
	err = h.checkTodayStages(ctx, todayWorkDate, offices)
	if err != nil {
		h.logger.Error("HourlyJob h.checkTodayStages error", zap.Error(err))
		return err
	}

	// чекаем что сегодня прошло 14:00 и открываем запись
	err = h.checkTomorrowStages(ctx, tomorrowWorkDate, offices)
	if err != nil {
		h.logger.Error("HourlyJob h.checkTomorrowStages error", zap.Error(err))
		return err
	}

	return nil
}

func (h *hourlyJobImpl) openTodayAccept(ctx context.Context, today model.WorkDate) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	if today.Status == model.StatusSetBookSeats {
		err := h.workDateRepo.UpdateStatusByID(today.ID, model.StatusAccept)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *hourlyJobImpl) checkTodayStages(ctx context.Context, today model.WorkDate, offices []*model.Office) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	for _, office := range offices {
		currentTime, err := helper.CurrentTimeWithTimeZone(office.TimeZone)
		morningNotify, err := helper.TimeWithTimeZone(helper.SendNotifyTime, office.TimeZone)
		removeBookTime, err := helper.TimeWithTimeZone(helper.RemoveBookTime, office.TimeZone)
		evening, err := helper.TimeWithTimeZone(helper.Evening, office.TimeZone)

		// Send notify users to confirm
		if currentTime.After(morningNotify) || currentTime.Equal(morningNotify) {
			err = h.informerService.SendNotifiesToConfirm(ctx, office)
			if err != nil {
				return err
			}
		}

		// Cancel book time
		if currentTime.After(removeBookTime) || currentTime.Equal(removeBookTime) {
			err = h.cancelNotConfirmedBookSeats(ctx, today, office)
			if err != nil {
				return err
			}
		}

		// Update to done status after 18-00 o'clock
		if currentTime.After(evening) || currentTime.Equal(evening) {
			err = h.workDateRepo.UpdateStatusByID(today.ID, model.StatusDone)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *hourlyJobImpl) checkTomorrowStages(ctx context.Context, tomorrow model.WorkDate, offices []*model.Office) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	for _, office := range offices {
		currentTime, err := helper.CurrentTimeWithTimeZone(office.TimeZone)
		openBooking, err := helper.TimeWithTimeZone(helper.OpenBooking, office.TimeZone)

		if (currentTime.After(openBooking) || currentTime.Equal(openBooking)) && tomorrow.Status != model.StatusAccept {
			err = h.workDateRepo.UpdateStatusByID(tomorrow.ID, model.StatusAccept)
			if err != nil {
				return err
			}

			formattedDate := tomorrow.Date.Format(helper.DateFormat)
			message := fmt.Sprintf("Открыта запись на %s в офис: %s", formattedDate, office.Name)
			err = h.informerService.SendNotifyTomorrowBookingOpen(ctx, *office, message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *hourlyJobImpl) cancelNotConfirmedBookSeats(ctx context.Context, today model.WorkDate, office *model.Office) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	if office == nil {
		return nil
	}

	bookSeats, err := h.bookSeatRepo.FindNotConfirmedByOfficeIDAndDate(ctx, office.ID, today.Date.String())
	if err != nil {
		return err
	}

	for _, bookSeat := range bookSeats {
		err = h.bookSeatRepo.CancelBookSeatWithID(ctx, bookSeat.ID)
		if err != nil {
			return err
		}
	}

	return h.informerService.SendNotifyToBookDeletedBySystem(ctx, bookSeats, office.Name)
}

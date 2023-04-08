package job

import (
	"fmt"
	"go.uber.org/zap"
	"telegram-api/internal/app/informer"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/common"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/repo/interfaces"
)

type hourlyJobImpl struct {
	informerService informer.InformerService
	officeRepo      interfaces.OfficeRepository
	workDateRepo    interfaces.WorkDateRepository
	logger          *zap.Logger
}

type HourlyJob interface {
	StartSchedule() error
}

func NewHourlyJob(informerService informer.InformerService,
	officeRepo interfaces.OfficeRepository,
	workDateRepo interfaces.WorkDateRepository,
	logger *zap.Logger) HourlyJob {
	return &hourlyJobImpl{
		informerService: informerService,
		officeRepo:      officeRepo,
		workDateRepo:    workDateRepo,
		logger:          logger,
	}
}

func (h *hourlyJobImpl) StartSchedule() error {

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
	err = h.checkTodayOpened(todayWorkDate)
	if err != nil {
		h.logger.Error("HourlyJob h.checkTodayOpened error", zap.Error(err))
		return err
	}

	offices, err := h.officeRepo.GetAll()
	if err != nil {
		h.logger.Error("HourlyJob officeRepo.GetAll error", zap.Error(err))
		return err
	}

	// чекаем что сегодня рабочий день закончился и закрываем его
	err = h.checkTodayStages(todayWorkDate, offices)
	if err != nil {
		h.logger.Error("HourlyJob h.checkTodayStages error", zap.Error(err))
		return err
	}

	// чекаем что сегодня прошло 14:00 и открываем запись
	err = h.checkTomorrowStages(tomorrowWorkDate, offices)
	if err != nil {
		h.logger.Error("HourlyJob h.checkTomorrowStages error", zap.Error(err))
		return err
	}

	return nil
}

func (h *hourlyJobImpl) checkTodayOpened(today model.WorkDate) error {
	if today.Status == model.StatusSetBookSeats {
		err := h.workDateRepo.UpdateStatusByID(today.ID, model.StatusAccept)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *hourlyJobImpl) checkTodayStages(today model.WorkDate, offices []*model.Office) error {
	for _, office := range offices {
		currentTime, err := helper.CurrentTimeWithTimeZone(office.TimeZone)
		evening, err := helper.TimeWithTimeZone(helper.Evening, office.TimeZone)

		if currentTime.After(evening) || currentTime.Equal(evening) {
			err = h.workDateRepo.UpdateStatusByID(today.ID, model.StatusDone)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *hourlyJobImpl) checkTomorrowStages(tomorrow model.WorkDate, offices []*model.Office) error {
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
			err = h.informerService.SendNotifiesWithMessage(*office, message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

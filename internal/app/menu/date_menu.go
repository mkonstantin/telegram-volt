package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type dateMenuImpl struct {
	dateRepo     repo.WorkDateRepository
	officeRepo   repo.OfficeRepository
	bookSeatRepo repo.BookSeatRepository
	dateMenuForm form.DateMenuForm
	cfg          config.AppConfig
	logger       *zap.Logger
}

func NewDateMenu(
	dateRepo repo.WorkDateRepository,
	officeRepo repo.OfficeRepository,
	bookSeatRepo repo.BookSeatRepository,
	dateMenuForm form.DateMenuForm,
	cfg config.AppConfig,
	logger *zap.Logger) interfaces.DateMenu {

	return &dateMenuImpl{
		dateRepo:     dateRepo,
		officeRepo:   officeRepo,
		bookSeatRepo: bookSeatRepo,
		dateMenuForm: dateMenuForm,
		cfg:          cfg,
		logger:       logger,
	}
}

func (f *dateMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {

	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	currentUser := model.GetCurrentUser(ctx)

	office, err := f.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	dates, err := f.getDates(ctx, f.cfg.IsAdmin(currentUser.TelegramName))

	var daySeats []form.DaySeat
	for _, date := range dates {
		seats, err := f.bookSeatRepo.GetFreeSeatsByOfficeIDAndDate(ctx, currentUser.OfficeID, date.Date.String())
		if err != nil {
			return nil, err
		}
		daySeat := form.DaySeat{
			Date:        date.Date,
			SeatsNumber: len(seats),
		}
		daySeats = append(daySeats, daySeat)
	}

	var message string
	if len(daySeats) > 0 {
		message = fmt.Sprintf("Выберите дату:")
	} else {
		message = fmt.Sprintf("В этом офисе сегодня мест нет или не работает")
	}

	formData := form.DateMenuFormData{
		Message:     message,
		Office:      office,
		SeatByDates: daySeats,
	}
	return f.dateMenuForm.Build(ctx, formData)
}

func (f *dateMenuImpl) getDates(ctx context.Context, isAdmin bool) ([]model.WorkDate, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	today := helper.TodayZeroTimeUTC()

	var dates []model.WorkDate
	var err error

	if isAdmin {
		todayPlus10 := helper.TodayPlusUTC(10)
		dates, err = f.dateRepo.FindByDates(today.String(), todayPlus10.String())
		if err != nil {
			return nil, err
		}
	} else {
		todayPlus2 := helper.TodayPlusUTC(2)
		dates, err = f.dateRepo.FindByDatesAndStatus(today.String(), todayPlus2.String(), model.StatusAccept)
		if err != nil {
			return nil, err
		}
	}

	return dates, nil
}

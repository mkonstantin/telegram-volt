package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
)

type dateMenuImpl struct {
	dateRepo     repo.WorkDateRepository
	officeRepo   repo.OfficeRepository
	bookSeatRepo repo.BookSeatRepository
	dateMenuForm form.DateMenuForm
	logger       *zap.Logger
}

func NewDateMenu(
	dateRepo repo.WorkDateRepository,
	officeRepo repo.OfficeRepository,
	bookSeatRepo repo.BookSeatRepository,
	dateMenuForm form.DateMenuForm,
	logger *zap.Logger) interfaces.DateMenu {

	return &dateMenuImpl{
		dateRepo:     dateRepo,
		officeRepo:   officeRepo,
		bookSeatRepo: bookSeatRepo,
		dateMenuForm: dateMenuForm,
		logger:       logger,
	}
}

func (f *dateMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {

	currentUser := model.GetCurrentUser(ctx)

	office, err := f.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	today := helper.TodayZeroTimeUTC()
	todayPlus2 := helper.TodayPlusUTC(2)

	// получаем сегодня и завтра
	dates, err := f.dateRepo.FindByDatesAndStatus(today.String(), todayPlus2.String(), model.StatusAccept)
	if err != nil {
		return nil, err
	}

	var seatByDates []form.DaySeat
	for _, date := range dates {
		seats, err := f.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, date.Date.String())
		if err != nil {
			return nil, err
		}
		daySeats := form.DaySeat{
			Date:        date.Date,
			SeatsNumber: len(seats),
		}
		seatByDates = append(seatByDates, daySeats)
	}

	var message string
	if len(seatByDates) > 0 {
		message = fmt.Sprintf("Выберите дату:")
	} else {
		message = fmt.Sprintf("В этом офисе сегодня мест нет или не работает")
	}

	formData := form.DateMenuFormData{
		Message:     message,
		Office:      office,
		SeatByDates: seatByDates,
	}
	return f.dateMenuForm.Build(ctx, formData)
}

package menu

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/internal/infrastructure/service"
)

type dateMenuImpl struct {
	officeRepo   repo.OfficeRepository
	bookSeatRepo repo.BookSeatRepository
	dateMenuForm form.DateMenuForm
	logger       *zap.Logger
}

func NewDateMenu(
	officeRepo repo.OfficeRepository,
	bookSeatRepo repo.BookSeatRepository,
	dateMenuForm form.DateMenuForm,
	logger *zap.Logger) interfaces.DateMenu {

	return &dateMenuImpl{
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

	today := service.TodayZeroTimeUTC()
	todaySeats, err := f.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, today.String())
	if err != nil {
		return nil, err
	}

	tomorrow := service.TomorrowZeroTimeUTC()
	tomorrowSeats, err := f.bookSeatRepo.GetAllByOfficeIDAndDate(currentUser.OfficeID, tomorrow.String())
	if err != nil {
		return nil, err
	}

	todayD := form.DaySeat{
		Date:        today.String(),
		SeatsNumber: len(todaySeats),
	}

	tomorrowD := form.DaySeat{
		Date:        tomorrow.String(),
		SeatsNumber: len(tomorrowSeats),
	}

	var seatByDates []form.DaySeat
	seatByDates = append(seatByDates, todayD)
	seatByDates = append(seatByDates, tomorrowD)

	message := fmt.Sprintf("Выберите дату:")

	formData := form.DateMenuFormData{
		Message:     message,
		Office:      office,
		SeatByDates: seatByDates,
	}
	return f.dateMenuForm.Build(ctx, formData)
}

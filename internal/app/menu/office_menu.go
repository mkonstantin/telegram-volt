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

type officeMenuImpl struct {
	dateRepo       repo.WorkDateRepository
	bookSeatRepo   repo.BookSeatRepository
	officeRepo     repo.OfficeRepository
	officeMenuForm form.OfficeMenuForm
	logger         *zap.Logger
}

func NewOfficeMenu(
	dateRepo repo.WorkDateRepository,
	bookSeatRepo repo.BookSeatRepository,
	officeRepo repo.OfficeRepository,
	officeMenuForm form.OfficeMenuForm,
	logger *zap.Logger) interfaces.OfficeMenu {

	return &officeMenuImpl{
		dateRepo:       dateRepo,
		bookSeatRepo:   bookSeatRepo,
		officeRepo:     officeRepo,
		officeMenuForm: officeMenuForm,
		logger:         logger,
	}
}

func (o *officeMenuImpl) Call(ctx context.Context, title string) (*tgbotapi.MessageConfig, error) {
	currentUser := model.GetCurrentUser(ctx)

	office, err := o.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	dates, err := o.dateRepo.FindByStatus(model.StatusAccept)
	if err != nil {
		return nil, err
	}

	today := helper.TodayZeroTimeUTC()

	var needConfirmBookSeat *model.BookSeat

	var bookSeats []*model.BookSeat
	for _, date := range dates {
		bookSeat, err := o.bookSeatRepo.FindByUserIDAndDate(currentUser.ID, date.Date.String())
		if err != nil {
			return nil, err
		}
		if bookSeat != nil {
			bookSeats = append(bookSeats, bookSeat)

			// Проверяем нужно ли подтверждение места на сегодня
			if bookSeat.BookDate == today && !bookSeat.Confirm {
				currentTime, err := helper.CurrentTimeWithTimeZone(bookSeat.Office.TimeZone)
				morningTime, err := helper.TimeWithTimeZone(helper.Morning, bookSeat.Office.TimeZone)
				if err != nil {
					return nil, err
				}
				if currentTime.After(morningTime) || currentTime.Equal(morningTime) {
					needConfirmBookSeat = bookSeat
				}
			}
		}
	}

	var buttonText string
	if currentUser.OfficeID == currentUser.NotifyOfficeID {
		buttonText = "Отписаться от уведомлений"
	} else {
		buttonText = "Подписаться на свободные места"
	}

	var message string
	if title != "" {
		message = title
	} else {
		message = fmt.Sprintf("Офис: %s, действия:", office.Name)
	}

	data := form.OfficeMenuFormData{
		Office:              office,
		Message:             message,
		SubscribeButtonText: buttonText,
		BookSeats:           bookSeats,
		NeedConfirmBookSeat: needConfirmBookSeat,
	}
	return o.officeMenuForm.Build(ctx, data)
}

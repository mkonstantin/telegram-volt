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

func (o *officeMenuImpl) Call(ctx context.Context) (*tgbotapi.MessageConfig, error) {
	currentUser := model.GetCurrentUser(ctx)

	office, err := o.officeRepo.FindByID(currentUser.OfficeID)
	if err != nil {
		return nil, err
	}

	dates, err := o.dateRepo.FindByStatus(model.StatusAccept)
	if err != nil {
		return nil, err
	}

	var bookSeats []*model.BookSeat
	for _, date := range dates {
		bookSeat, err := o.bookSeatRepo.FindByUserIDAndDate(currentUser.ID, date.Date.String())
		if err != nil {
			return nil, err
		}
		if bookSeat != nil {
			bookSeats = append(bookSeats, bookSeat)
		}
	}

	var buttonText string
	if currentUser.OfficeID == currentUser.NotifyOfficeID {
		buttonText = "Отписаться от уведомлений"
	} else {
		buttonText = "Подписаться на свободные места"
	}

	message := fmt.Sprintf("Офис: %s, действия:", office.Name)

	data := form.OfficeMenuFormData{
		Office:              office,
		Message:             message,
		SubscribeButtonText: buttonText,
		BookSeats:           bookSeats,
	}
	return o.officeMenuForm.Build(ctx, data)
}

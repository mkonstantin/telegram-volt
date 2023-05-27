package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/informer"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
	"telegram-api/pkg/tracing"
)

type OfficeMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeMenuImpl struct {
	informerService informer.InformerService
	userService     usecase.UserService
	dateMenu        interfaces.DateMenu
	officeMenu      interfaces.OfficeMenu
	officeListMenu  interfaces.OfficeListMenu
	logger          *zap.Logger
}

func NewOfficeMenuHandle(
	informerService informer.InformerService,
	userService usecase.UserService,
	dateMenu interfaces.DateMenu,
	officeMenu interfaces.OfficeMenu,
	officeListMenu interfaces.OfficeListMenu,
	logger *zap.Logger) OfficeMenu {

	return &officeMenuImpl{
		informerService: informerService,
		userService:     userService,
		dateMenu:        dateMenu,
		officeMenu:      officeMenu,
		officeListMenu:  officeListMenu,
		logger:          logger,
	}
}

func (o *officeMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	switch command.Action {
	case dto.OfficeMenuFreeSeats:
		return o.dateMenu.Call(ctx)

	case dto.OfficeMenuSubscribe:
		contxt, message, err := o.userService.SubscribeWork(ctx)
		if err != nil {
			return nil, err
		}
		return o.officeMenu.Call(contxt, message, 0)

	case dto.OfficeMenuChooseAnotherOffice:
		return o.officeListMenu.Call(ctx)

	case dto.OfficeMenuCancelBook:
		message, isCanceled, err := o.userService.CancelBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		if isCanceled {
			err = o.informerService.SendNotifySeatBecomeFree(ctx, command.BookSeatID)
			if err != nil {
				return nil, err
			}
			return o.officeMenu.Call(ctx, message, 0)
		} else {
			chatID := model.GetCurrentChatID(ctx)
			msg := tgbotapi.NewMessage(chatID, "")
			msg.Text = message
			return &msg, nil
		}

	case dto.OfficeMenuConfirm:
		message, err := o.userService.ConfirmBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}

		return o.officeMenu.Call(ctx, message, 0)
	}

	return nil, nil
}

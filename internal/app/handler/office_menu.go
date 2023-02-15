package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
)

type OfficeMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeMenuImpl struct {
	userService    usecase.UserService
	seatListMenu   interfaces.SeatListMenu
	officeListMenu interfaces.OfficeListMenu
	logger         *zap.Logger
}

func NewOfficeMenuHandle(
	userService usecase.UserService,
	seatListMenu interfaces.SeatListMenu,
	officeListMenu interfaces.OfficeListMenu,
	logger *zap.Logger) OfficeMenu {

	return &officeMenuImpl{
		userService:    userService,
		seatListMenu:   seatListMenu,
		officeListMenu: officeListMenu,
		logger:         logger,
	}
}

func (o *officeMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.OfficeMenuFreeSeats:
		return o.seatListMenu.Call(ctx)

	case dto.OfficeMenuSubscribe:
		message, err := o.userService.SubscribeWork(ctx)
		if err != nil {
			return nil, err
		}

		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil

	case dto.OfficeMenuChooseAnotherOffice:
		return o.officeListMenu.Call(ctx)
	}

	return nil, nil
}

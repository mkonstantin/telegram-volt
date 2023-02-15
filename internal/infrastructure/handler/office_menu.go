package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/former"
	"telegram-api/internal/infrastructure/handler/dto"
)

type OfficeMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeMenuImpl struct {
	userService    usecase.UserService
	officeListMenu interfaces.OfficeListMenu
	msgFormer      former.MessageFormer
	logger         *zap.Logger
}

func NewOfficeMenuHandle(
	userService usecase.UserService,
	officeListMenu interfaces.OfficeListMenu,
	msgFormer former.MessageFormer,
	logger *zap.Logger) OfficeMenu {

	return &officeMenuImpl{
		userService:    userService,
		officeListMenu: officeListMenu,
		msgFormer:      msgFormer,
		logger:         logger,
	}
}

func (o *officeMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.OfficeMenuFreeSeats:
		result, err := o.userService.CallDateMenu(ctx)
		if err != nil {
			return nil, err
		}
		return o.msgFormer.FormSeatListMsg(ctx, result)

	case dto.OfficeMenuSubscribe:
		result, err := o.userService.SubscribeWork(ctx)
		if err != nil {
			return nil, err
		}

		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = result.Message

		return &msg, nil

	case dto.OfficeMenuChooseAnotherOffice:
		return o.officeListMenu.Call(ctx)
	}

	return nil, nil
}

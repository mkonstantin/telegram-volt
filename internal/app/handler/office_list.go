package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
)

type OfficeList interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeListImpl struct {
	userService usecase.UserService
	officeMenu  interfaces.OfficeMenu
	logger      *zap.Logger
}

func NewOfficeListHandle(
	userService usecase.UserService,
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) OfficeList {

	return &officeListImpl{
		userService: userService,
		officeMenu:  officeMenu,
		logger:      logger,
	}
}

func (o *officeListImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	ctx, err := o.userService.SetOfficeScript(ctx, command.OfficeID)
	if err != nil {
		return nil, err
	}

	return o.officeMenu.Call(ctx, "", 0)
}

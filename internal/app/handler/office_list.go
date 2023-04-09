package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
)

type OfficeList interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type officeListImpl struct {
	officeMenu interfaces.OfficeMenu
	logger     *zap.Logger
}

func NewOfficeListHandle(
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) OfficeList {

	return &officeListImpl{
		officeMenu: officeMenu,
		logger:     logger,
	}
}

func (o *officeListImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	return o.officeMenu.Call(ctx, "", command.OfficeID)
}

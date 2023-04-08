package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
)

type DateMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type dateMenuImpl struct {
	seatList   interfaces.SeatListMenu
	officeMenu interfaces.OfficeMenu
	logger     *zap.Logger
}

func NewDateMenuHandle(
	seatList interfaces.SeatListMenu,
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) DateMenu {

	return &dateMenuImpl{
		seatList:   seatList,
		officeMenu: officeMenu,
		logger:     logger,
	}
}

func (o *dateMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	if command.Action == dto.Back {
		return o.officeMenu.Call(ctx, "")
	}

	if command.BookDate == nil {
		return o.officeMenu.Call(ctx, "")
	}

	return o.seatList.Call(ctx, *command.BookDate)
}

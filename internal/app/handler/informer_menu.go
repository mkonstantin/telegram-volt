package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
)

type InfoMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type infoMenuImpl struct {
	seatList   interfaces.SeatListMenu
	officeMenu interfaces.OfficeMenu
	logger     *zap.Logger
}

func NewInfoMenuHandle(
	seatList interfaces.SeatListMenu,
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) InfoMenu {

	return &infoMenuImpl{
		seatList:   seatList,
		officeMenu: officeMenu,
		logger:     logger,
	}
}

func (o *infoMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	if command.BookDate == nil {
		return o.officeMenu.Call(ctx)
	}

	return o.seatList.Call(ctx, *command.BookDate)
}

package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/infrastructure/former"
	"telegram-api/internal/infrastructure/handler/dto"
)

type FreeSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type freeSeatMenuImpl struct {
	userService  usecase.UserService
	seatListMenu interfaces.SeatListMenu
	msgFormer    former.MessageFormer
	logger       *zap.Logger
}

func NewFreeSeatMenuHandle(
	userService usecase.UserService,
	seatListMenu interfaces.SeatListMenu,
	msgFormer former.MessageFormer,
	logger *zap.Logger) FreeSeatMenu {

	return &freeSeatMenuImpl{
		userService:  userService,
		seatListMenu: seatListMenu,
		msgFormer:    msgFormer,
		logger:       logger,
	}
}

func (f *freeSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionBookYes:
		result, err := f.userService.BookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return f.msgFormer.FormBookSeatResult(ctx, result)
	case dto.ActionBookNo:
		fallthrough
	default:
		return f.seatListMenu.Call(ctx)
	}
}

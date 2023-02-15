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

type OwnSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type ownSeatMenuImpl struct {
	userService  usecase.UserService
	seatListMenu interfaces.SeatListMenu
	msgFormer    former.MessageFormer
	logger       *zap.Logger
}

func NewOwnSeatMenuHandle(
	userService usecase.UserService,
	seatListMenu interfaces.SeatListMenu,
	msgFormer former.MessageFormer,
	logger *zap.Logger) OwnSeatMenu {

	return &ownSeatMenuImpl{
		userService:  userService,
		seatListMenu: seatListMenu,
		msgFormer:    msgFormer,
		logger:       logger,
	}
}

func (o *ownSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionCancelBookYes:
		result, err := o.userService.CancelBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return o.msgFormer.FormCancelBookResult(ctx, result)
	case dto.ActionCancelBookNo:
		fallthrough
	default:
		return o.seatListMenu.Call(ctx)
	}
}

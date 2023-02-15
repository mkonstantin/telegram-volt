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

type OwnSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type ownSeatMenuImpl struct {
	userService  usecase.UserService
	seatListMenu interfaces.SeatListMenu
	logger       *zap.Logger
}

func NewOwnSeatMenuHandle(
	userService usecase.UserService,
	seatListMenu interfaces.SeatListMenu,
	logger *zap.Logger) OwnSeatMenu {

	return &ownSeatMenuImpl{
		userService:  userService,
		seatListMenu: seatListMenu,
		logger:       logger,
	}
}

func (o *ownSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionCancelBookYes:
		message, err := o.userService.CancelBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}

		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil

	case dto.ActionCancelBookNo:
		fallthrough
	default:
		return o.seatListMenu.Call(ctx)
	}
}

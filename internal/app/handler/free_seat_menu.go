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

type FreeSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type freeSeatMenuImpl struct {
	userService  usecase.UserService
	seatListMenu interfaces.SeatListMenu
	logger       *zap.Logger
}

func NewFreeSeatMenuHandle(
	userService usecase.UserService,
	seatListMenu interfaces.SeatListMenu,
	logger *zap.Logger) FreeSeatMenu {

	return &freeSeatMenuImpl{
		userService:  userService,
		seatListMenu: seatListMenu,
		logger:       logger,
	}
}

func (f *freeSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionBookYes:
		message, err := f.userService.BookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}

		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil

	case dto.ActionBookNo:
		fallthrough
	default:
		return f.seatListMenu.Call(ctx)
	}
}

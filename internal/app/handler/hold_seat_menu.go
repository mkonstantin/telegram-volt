package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
	interfaces2 "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type HoldSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type holdSeatMenuImpl struct {
	userService  usecase.UserService
	bookSeatRepo interfaces2.BookSeatRepository
	seatListMenu interfaces.SeatListMenu
	logger       *zap.Logger
}

func NewHoldSeatMenuHandle(
	userService usecase.UserService,
	bookSeatRepo interfaces2.BookSeatRepository,
	seatListMenu interfaces.SeatListMenu,
	logger *zap.Logger) HoldSeatMenu {

	return &holdSeatMenuImpl{
		userService:  userService,
		bookSeatRepo: bookSeatRepo,
		seatListMenu: seatListMenu,
		logger:       logger,
	}
}

func (o *holdSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	switch command.Action {
	case dto.ActionCancelHoldYes:
		message, err := o.userService.CancelHoldBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, message)
		return &msg, nil

	case dto.ActionCancelHoldNo:
		fallthrough

	default:
		bookSeat, err := o.bookSeatRepo.FindByID(command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return o.seatListMenu.Call(ctx, bookSeat.BookDate, 0)
	}
}

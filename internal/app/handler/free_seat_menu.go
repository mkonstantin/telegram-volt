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

type FreeSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type freeSeatMenuImpl struct {
	officeMenu   interfaces.OfficeMenu
	userService  usecase.UserService
	bookSeatRepo interfaces2.BookSeatRepository
	seatListMenu interfaces.SeatListMenu
	logger       *zap.Logger
}

func NewFreeSeatMenuHandle(
	officeMenu interfaces.OfficeMenu,
	userService usecase.UserService,
	bookSeatRepo interfaces2.BookSeatRepository,
	seatListMenu interfaces.SeatListMenu,
	logger *zap.Logger) FreeSeatMenu {

	return &freeSeatMenuImpl{
		officeMenu:   officeMenu,
		userService:  userService,
		bookSeatRepo: bookSeatRepo,
		seatListMenu: seatListMenu,
		logger:       logger,
	}
}

func (f *freeSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	switch command.Action {
	case dto.ActionBookYes:
		message, err := f.userService.BookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return f.officeMenu.Call(ctx, message, 0)

	case dto.ActionBookHold:
		message, err := f.userService.HoldBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, message)
		return &msg, nil

	case dto.ActionBookNo:
		fallthrough
	default:
		bookSeat, err := f.bookSeatRepo.FindByID(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return f.seatListMenu.Call(ctx, bookSeat.BookDate, 0)
	}
}

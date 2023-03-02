package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/service"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
	interfaces2 "telegram-api/internal/infrastructure/repo/interfaces"
)

type OwnSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type ownSeatMenuImpl struct {
	informerService service.InformerService
	userService     usecase.UserService
	bookSeatRepo    interfaces2.BookSeatRepository
	seatListMenu    interfaces.SeatListMenu
	logger          *zap.Logger
}

func NewOwnSeatMenuHandle(
	informerService service.InformerService,
	userService usecase.UserService,
	bookSeatRepo interfaces2.BookSeatRepository,
	seatListMenu interfaces.SeatListMenu,
	logger *zap.Logger) OwnSeatMenu {

	return &ownSeatMenuImpl{
		informerService: informerService,
		userService:     userService,
		bookSeatRepo:    bookSeatRepo,
		seatListMenu:    seatListMenu,
		logger:          logger,
	}
}

func (o *ownSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	switch command.Action {
	case dto.ActionCancelBookYes:
		message, isCanceled, err := o.userService.CancelBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		if isCanceled {
			err = o.informerService.SeatComeFree(ctx, command.BookSeatID)
			if err != nil {
				return nil, err
			}
		}
		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil

	case dto.ActionCancelBookNo:
		fallthrough
	default:
		bookSeat, err := o.bookSeatRepo.FindByID(command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return o.seatListMenu.Call(ctx, bookSeat.BookDate)
	}
}

package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/informer"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/app/usecase"
	interfaces2 "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type OwnSeatMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type ownSeatMenuImpl struct {
	officeMenu      interfaces.OfficeMenu
	informerService informer.InformerService
	userService     usecase.UserService
	bookSeatRepo    interfaces2.BookSeatRepository
	seatListMenu    interfaces.SeatListMenu
	logger          *zap.Logger
}

func NewOwnSeatMenuHandle(
	officeMenu interfaces.OfficeMenu,
	informerService informer.InformerService,
	userService usecase.UserService,
	bookSeatRepo interfaces2.BookSeatRepository,
	seatListMenu interfaces.SeatListMenu,
	logger *zap.Logger) OwnSeatMenu {

	return &ownSeatMenuImpl{
		officeMenu:      officeMenu,
		informerService: informerService,
		userService:     userService,
		bookSeatRepo:    bookSeatRepo,
		seatListMenu:    seatListMenu,
		logger:          logger,
	}
}

func (o *ownSeatMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	switch command.Action {
	case dto.ActionCancelBookYes:
		message, isCanceled, err := o.userService.CancelBookSeat(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		if isCanceled {
			err = o.informerService.SendNotifySeatBecomeFree(ctx, command.BookSeatID)
			if err != nil {
				return nil, err
			}
		}
		return o.officeMenu.Call(ctx, message, 0)

	case dto.ActionCancelBookNo:
		fallthrough
	default:
		bookSeat, err := o.bookSeatRepo.FindByID(ctx, command.BookSeatID)
		if err != nil {
			return nil, err
		}
		return o.seatListMenu.Call(ctx, bookSeat.BookDate, 0)
	}
}

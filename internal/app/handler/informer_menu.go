package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	interfaces2 "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

type InfoMenu interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type infoMenuImpl struct {
	bookSeatRepo interfaces2.BookSeatRepository
	seatList     interfaces.SeatListMenu
	officeMenu   interfaces.OfficeMenu
	logger       *zap.Logger
}

func NewInfoMenuHandle(
	bookSeatRepo interfaces2.BookSeatRepository,
	seatList interfaces.SeatListMenu,
	officeMenu interfaces.OfficeMenu,
	logger *zap.Logger) InfoMenu {

	return &infoMenuImpl{
		bookSeatRepo: bookSeatRepo,
		seatList:     seatList,
		officeMenu:   officeMenu,
		logger:       logger,
	}
}

func (o *infoMenuImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	bookSeatID := command.BookSeatID

	if bookSeatID == 0 {
		currentUser := model.GetCurrentUser(ctx)
		return o.officeMenu.Call(ctx, "", currentUser.NotifyOfficeID)
	}

	bookSeat, err := o.bookSeatRepo.FindByID(ctx, bookSeatID)
	if err != nil {
		return nil, err
	}

	switch command.Action {
	case dto.ActionShowOfficeMenu:
		return o.officeMenu.Call(ctx, "", bookSeat.Office.ID)
	case dto.ActionShowSeatList:
		fallthrough
	default:
		return o.seatList.Call(ctx, bookSeat.BookDate, bookSeat.Office.ID)
	}
}

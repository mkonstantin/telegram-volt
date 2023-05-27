package handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
	"telegram-api/pkg/tracing"
)

const (
	ThisIsYourSeat = "this_is_your_seat"
	ThisIsSeatBusy = "this_is_seat_busy"
	ThisIsSeatFree = "this_is_seat_free"
	ThisIsSeatHold = "this_is_seat_hold"
)

type SeatList interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type seatListImpl struct {
	bookSeatRepo repo.BookSeatRepository
	dateMenu     interfaces.DateMenu
	ownSeatMenu  interfaces.OwnSeatMenu
	holdSeatMenu interfaces.HoldSeatMenu
	freeSeatMenu interfaces.FreeSeatMenu
	cfg          config.AppConfig
	logger       *zap.Logger
}

func NewSeatListHandle(
	bookSeatRepo repo.BookSeatRepository,
	dateMenu interfaces.DateMenu,
	ownSeatMenu interfaces.OwnSeatMenu,
	holdSeatMenu interfaces.HoldSeatMenu,
	freeSeatMenu interfaces.FreeSeatMenu,
	cfg config.AppConfig,
	logger *zap.Logger) SeatList {

	return &seatListImpl{
		bookSeatRepo: bookSeatRepo,
		dateMenu:     dateMenu,
		ownSeatMenu:  ownSeatMenu,
		holdSeatMenu: holdSeatMenu,
		freeSeatMenu: freeSeatMenu,
		cfg:          cfg,
		logger:       logger,
	}
}

func (s *seatListImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	if command.Action == dto.Back {
		return s.dateMenu.Call(ctx)
	}

	bookSeat, err := s.bookSeatRepo.FindByID(ctx, command.BookSeatID)
	if err != nil {
		return nil, err
	}

	chatID := model.GetCurrentChatID(ctx)

	switch getStatus(ctx, bookSeat) {
	case ThisIsSeatFree:
		return s.freeSeatMenu.Call(ctx, command.BookSeatID)
	case ThisIsYourSeat:
		return s.ownSeatMenu.Call(ctx, command.BookSeatID)
	case ThisIsSeatHold:
		return s.seatHold(ctx, bookSeat)
	case ThisIsSeatBusy:
		fallthrough
	default:
		message := fmt.Sprintf("Место №%s уже занято %s aka @%s",
			bookSeat.Seat.SeatSign, bookSeat.User.Name, bookSeat.User.TelegramName)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil
	}
}

func getStatus(ctx context.Context, bookSeat *model.BookSeat) string {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	currentUser := model.GetCurrentUser(ctx)

	if bookSeat.User != nil {
		if bookSeat.User.TelegramID == currentUser.TelegramID {
			return ThisIsYourSeat
		} else {
			return ThisIsSeatBusy
		}
	} else {
		if bookSeat.Hold {
			return ThisIsSeatHold
		}
		return ThisIsSeatFree
	}
}

func (s *seatListImpl) seatHold(ctx context.Context, bookSeat *model.BookSeat) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	currentUser := model.GetCurrentUser(ctx)

	if s.cfg.IsAdmin(currentUser.TelegramName) {
		return s.holdSeatMenu.Call(ctx, bookSeat.ID)
	} else {
		chatID := model.GetCurrentChatID(ctx)
		message := fmt.Sprintf("Место №%s временно закреплено администратором, его нельзя забронировать", bookSeat.Seat.SeatSign)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil
	}
}

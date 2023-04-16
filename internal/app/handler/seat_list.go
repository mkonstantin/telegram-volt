package handler

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/app/menu/interfaces"
	"telegram-api/internal/domain/model"
	repo "telegram-api/internal/infrastructure/repo/interfaces"
)

const (
	ThisIsYourSeat = "this_is_your_seat"
	ThisIsSeatBusy = "this_is_seat_busy"
	ThisIsSeatFree = "this_is_seat_free"
)

type SeatList interface {
	Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error)
}

type seatListImpl struct {
	bookSeatRepo repo.BookSeatRepository
	dateMenu     interfaces.DateMenu
	ownSeatMenu  interfaces.OwnSeatMenu
	freeSeatMenu interfaces.FreeSeatMenu
	logger       *zap.Logger
}

func NewSeatListHandle(
	bookSeatRepo repo.BookSeatRepository,
	dateMenu interfaces.DateMenu,
	ownSeatMenu interfaces.OwnSeatMenu,
	freeSeatMenu interfaces.FreeSeatMenu,
	logger *zap.Logger) SeatList {

	return &seatListImpl{
		bookSeatRepo: bookSeatRepo,
		dateMenu:     dateMenu,
		ownSeatMenu:  ownSeatMenu,
		freeSeatMenu: freeSeatMenu,
		logger:       logger,
	}
}

func (s *seatListImpl) Handle(ctx context.Context, command dto.InlineRequest) (*tgbotapi.MessageConfig, error) {

	if command.Action == dto.Back {
		return s.dateMenu.Call(ctx)
	}

	bookSeat, err := s.bookSeatRepo.FindByID(command.BookSeatID)
	if err != nil {
		return nil, err
	}

	switch getStatus(ctx, bookSeat) {
	case ThisIsSeatFree:
		return s.freeSeatMenu.Call(ctx, command.BookSeatID)
	case ThisIsYourSeat:
		return s.ownSeatMenu.Call(ctx, command.BookSeatID)
	case ThisIsSeatBusy:
		fallthrough
	default:
		message := fmt.Sprintf("Место №%s уже занято %s aka @%s",
			bookSeat.Seat.SeatSign, bookSeat.User.Name, bookSeat.User.TelegramName)

		chatID := model.GetCurrentChatID(ctx)
		msg := tgbotapi.NewMessage(chatID, "")
		msg.Text = message
		return &msg, nil
	}
}

func getStatus(ctx context.Context, bookSeat *model.BookSeat) string {
	currentUser := model.GetCurrentUser(ctx)

	if bookSeat.User != nil {
		if bookSeat.User.TelegramID == currentUser.TelegramID {
			return ThisIsYourSeat
		} else {
			return ThisIsSeatBusy
		}
	} else {
		return ThisIsSeatFree
	}
}

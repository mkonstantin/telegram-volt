package dto

import "telegram-api/internal/domain/model"

type FirstStartDTO struct {
	User      model.User
	MessageID int
	ChatID    int64
}

type BookSeatTapDTO struct {
	TelegramID int64
	BookSeatID int64
	ChatID     int64
	MessageID  int
}

type BookSeatDTO struct {
	TelegramID int64
	OfficeID   int64
	ChatID     int64
	MessageID  int
}

type OfficeChosenDTO struct {
	TelegramID int64
	OfficeID   int64
	ChatID     int64
	MessageID  int
}

type UserResult struct {
	Key        string
	Office     *model.Office
	Offices    []*model.Office
	BookSeats  []*model.BookSeat
	BookSeatID int64
	Message    string
	MessageID  int
	ChatID     int64
}

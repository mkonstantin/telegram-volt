package dto

import "telegram-api/internal/domain/model"

type FirstStartDTO struct {
	User      model.User
	MessageID int
	ChatID    int64
}

type FirstStartResult struct {
	Key       string
	Office    *model.Office
	Offices   []*model.Office
	Message   string
	User      model.User
	MessageID int
	ChatID    int64
}

type OfficeChosenDTO struct {
	TelegramID int64
	OfficeID   int64
}

type OfficeChosenResult struct {
	Seats []*model.Seat
}

package dto

import "telegram-api/internal/domain/model"

type FirstStartDTO struct {
	User      model.User
	MessageID int
	ChatID    int64
}

type OfficeChosenDTO struct {
	TelegramID int64
	OfficeID   int64
	ChatID     int64
	MessageID  int
}

type UserResult struct {
	Key       string
	Office    *model.Office
	Offices   []*model.Office
	Message   string
	MessageID int
	ChatID    int64
}

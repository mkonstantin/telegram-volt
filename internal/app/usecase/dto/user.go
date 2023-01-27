package dto

import "telegram-api/internal/domain/model"

type UserLogicData struct {
	User      model.User
	MessageID int
	ChatID    int64
}

type UserLogicResult struct {
	Key       string
	Office    *model.Office
	Offices   []*model.Office
	Message   string
	User      model.User
	MessageID int
	ChatID    int64
}

type SetOfficeDTO struct {
	TelegramID int64
	OfficeID   int64
}

package dto

import "telegram-api/internal/domain/model"

type UserResult struct {
	Key        string
	Office     *model.Office
	Offices    []*model.Office
	BookSeats  []*model.BookSeat
	BookSeatID int64
	Message    string
}

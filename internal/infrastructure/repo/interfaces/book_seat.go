package interfaces

import "telegram-api/internal/domain/model"

type BookSeatRepository interface {
	FindByID(id int64) (*model.BookSeat, error)
	GetAllByOfficeID(id int64) ([]*model.BookSeat, error)
	BookSeatWithID(id int64) (*model.BookSeat, error)
	CancelBookSeatWithID(id int64) (*model.BookSeat, error)
}
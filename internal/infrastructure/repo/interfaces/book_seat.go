package interfaces

import (
	"telegram-api/internal/domain/model"
	"time"
)

type BookSeatRepository interface {
	FindByID(id int64) (*model.BookSeat, error)
	GetAllByOfficeID(id int64) ([]*model.BookSeat, error)
	BookSeatWithID(id, userID int64) error
	CancelBookSeatWithID(id int64) error
	FindByUserID(userID int64) (*model.BookSeat, error)
	InsertSeat(officeID, seatID int64, dayDate time.Time) error
}

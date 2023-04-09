package interfaces

import (
	"telegram-api/internal/domain/model"
	"time"
)

type BookSeatRepository interface {
	FindByID(id int64) (*model.BookSeat, error)
	GetAllByOfficeIDAndDate(id int64, dateStr string) ([]*model.BookSeat, error)
	FindByOfficeIDAndDate(id int64, dateStr string) ([]*model.BookSeat, error)
	GetFreeSeatsByOfficeIDAndDate(id int64, dateStr string) ([]*model.BookSeat, error)
	BookSeatWithID(id, userID int64, confirm bool) error
	CancelBookSeatWithID(id int64) error
	FindByUserID(userID int64) (*model.BookSeat, error)
	FindByUserIDAndDate(userID int64, dateStr string) (*model.BookSeat, error)
	InsertSeat(officeID, seatID int64, dayDate time.Time) error
	ConfirmBookSeat(seatID int64) error
	FindNotConfirmedByOfficeIDAndDate(id int64, dateStr string) ([]*model.BookSeat, error)
}

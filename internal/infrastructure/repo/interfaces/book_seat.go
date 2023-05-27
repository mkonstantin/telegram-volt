package interfaces

import (
	"context"
	"telegram-api/internal/domain/model"
	"time"
)

type BookSeatRepository interface {
	FindByID(ctx context.Context, id int64) (*model.BookSeat, error)
	GetAllByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error)
	FindByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error)
	GetFreeSeatsByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error)
	BookSeatWithID(ctx context.Context, id, userID int64, confirm bool) error
	CancelBookSeatWithID(ctx context.Context, id int64) error
	FindByUserID(ctx context.Context, userID int64) (*model.BookSeat, error)
	FindByUserIDAndDate(ctx context.Context, userID int64, dateStr string) (*model.BookSeat, error)
	InsertSeat(ctx context.Context, officeID, seatID int64, dayDate time.Time) error
	ConfirmBookSeat(ctx context.Context, seatID int64) error
	FindNotConfirmedByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error)
	HoldSeatWithID(ctx context.Context, id int64) error
	CancelHoldSeatWithID(ctx context.Context, id int64) error
}

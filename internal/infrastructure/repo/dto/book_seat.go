package dto

import (
	"telegram-api/internal/domain/model"
	"time"
)

type BookSeat struct {
	ID               int64      `db:"id,omitempty"`
	OfficeID         int64      `db:"office_id,omitempty"`
	SeatID           int64      `db:"seat_id,omitempty"`
	BookDate         time.Time  `db:"book_date,omitempty"`
	SeatNumber       int        `db:"seat_number,omitempty"`
	IsHaveMonitor    bool       `db:"have_monitor,omitempty"`
	OfficeName       string     `db:"office_name,omitempty"`
	UserID           *int64     `db:"user_id,omitempty"`
	TelegramName     *string    `db:"user_name,omitempty"`
	TelegramID       *int64     `db:"telegram_id,omitempty"`
	TelegramUsername *string    `db:"telegram_name,omitempty"`
	BookStartTime    *time.Time `db:"book_start_time,omitempty"`
	BookEndTime      *time.Time `db:"book_end_time,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

func (o *BookSeat) ToModel() *model.BookSeat {
	var user *model.User
	if o.UserID != nil && o.TelegramName != nil && o.TelegramID != nil && o.TelegramUsername != nil {
		user = &model.User{
			ID:           *o.UserID,
			Name:         *o.TelegramName,
			TelegramID:   *o.TelegramID,
			TelegramName: *o.TelegramUsername,
		}
	}

	return &model.BookSeat{
		ID: o.ID,
		Office: model.Office{
			ID:   o.OfficeID,
			Name: o.OfficeName,
		},
		Seat: model.Seat{
			ID:            o.SeatID,
			IsHaveMonitor: o.IsHaveMonitor,
			SeatNumber:    o.SeatNumber,
			OfficeID:      o.OfficeID,
		},
		User:          user,
		BookDate:      o.BookDate,
		BookStartTime: o.BookStartTime,
		BookEndTime:   o.BookEndTime,
	}
}

func ToBookSeatModels(array []BookSeat) []*model.BookSeat {
	var models []*model.BookSeat
	for _, item := range array {
		models = append(models, item.ToModel())
	}
	return models
}

package dto

import (
	"telegram-api/internal/domain/model"
	"time"
)

type BookSeat struct {
	ID            int64     `db:"id,omitempty"`
	IsHaveMonitor bool      `db:"have_monitor,omitempty"`
	SeatNumber    int       `db:"seat_number,omitempty"`
	OfficeID      int64     `db:"office_id,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (o *BookSeat) ToModel() *model.BookSeat {
	return &model.BookSeat{
		ID:            0,
		Office:        model.Office{},
		Seat:          model.Seat{},
		User:          model.User{},
		BookDate:      time.Time{},
		BookStartTime: time.Time{},
		BookEndTime:   time.Time{},
	}
}

func ToBookSeatModels(array []BookSeat) []*model.BookSeat {
	var models []*model.BookSeat
	for _, item := range array {
		models = append(models, item.ToModel())
	}
	return models
}

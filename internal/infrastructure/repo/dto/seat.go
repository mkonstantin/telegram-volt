package dto

import (
	"telegram-api/internal/domain/model"
	"time"
)

type Seat struct {
	ID            int64     `db:"id,omitempty"`
	IsHaveMonitor bool      `db:"have_monitor,omitempty"`
	SeatSign      string    `db:"seat_sign,omitempty"`
	OfficeID      int64     `db:"office_id,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (o *Seat) ToModel() *model.Seat {
	return &model.Seat{
		ID:            o.ID,
		IsHaveMonitor: o.IsHaveMonitor,
		SeatSign:      o.SeatSign,
		OfficeID:      o.OfficeID,
	}
}

func ToSeatModels(array []Seat) []*model.Seat {
	var models []*model.Seat
	for _, item := range array {
		models = append(models, item.ToModel())
	}
	return models
}

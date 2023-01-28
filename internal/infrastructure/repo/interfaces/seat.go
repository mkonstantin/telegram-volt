package interfaces

import "telegram-api/internal/domain/model"

type SeatRepository interface {
	FindByID(id int64) (*model.Seat, error)
	GetAllByOfficeID(id int64) ([]*model.Seat, error)
}

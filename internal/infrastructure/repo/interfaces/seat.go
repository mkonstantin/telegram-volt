package interfaces

import (
	"context"
	"telegram-api/internal/domain/model"
)

type SeatRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Seat, error)
	GetAllByOfficeID(ctx context.Context, id int64) ([]*model.Seat, error)
}

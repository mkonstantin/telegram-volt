package interfaces

import (
	"context"
	"telegram-api/internal/domain/model"
	"time"
)

type WorkDateRepository interface {
	GetLastByDate(ctx context.Context) (*model.WorkDate, error)
	DoneAllPastByDate(ctx context.Context, date string) error
	FindByDates(ctx context.Context, startDate string, endDate string) ([]model.WorkDate, error)
	FindByDatesAndStatus(ctx context.Context, startDate string, endDate string, status model.DateStatus) ([]model.WorkDate, error)
	FindByStatus(ctx context.Context, status model.DateStatus) ([]model.WorkDate, error)
	InsertDate(ctx context.Context, dayDate time.Time) error
	FindByID(ctx context.Context, id int64) (*model.WorkDate, error)
	UpdateStatusByID(ctx context.Context, id int64, status model.DateStatus) error
}

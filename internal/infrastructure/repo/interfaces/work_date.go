package interfaces

import (
	"telegram-api/internal/domain/model"
	"time"
)

type WorkDateRepository interface {
	GetLastByDate() (*model.WorkDate, error)
	FindByDates(startDate string, endDate string) ([]model.WorkDate, error)
	FindByDatesAndStatus(startDate string, endDate string, status model.DateStatus) ([]model.WorkDate, error)
	InsertDate(dayDate time.Time) error
	FindByID(id int64) (*model.WorkDate, error)
	UpdateStatusByID(id int64, status model.DateStatus) error
}

package interfaces

import (
	"telegram-api/internal/domain/model"
	"time"
)

type WorkDateRepository interface {
	GetLastByDate() (*model.WorkDate, error)
	FindByDateAndStatus(dateStr string, status model.DateStatus, limit uint64) ([]model.WorkDate, error)
	InsertDate(dayDate time.Time) error
	FindByID(id int64) (*model.WorkDate, error)
}

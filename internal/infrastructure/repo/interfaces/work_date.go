package interfaces

import (
	"telegram-api/internal/domain/model"
	"time"
)

type WorkDateRepository interface {
	GetLastByDate() (*model.WorkDate, error)
	InsertDate(dayDate time.Time) error
	FindByID(id int64) (*model.WorkDate, error)
}

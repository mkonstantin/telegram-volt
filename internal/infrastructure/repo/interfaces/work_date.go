package interfaces

import "telegram-api/internal/domain/model"

type WorkDateRepository interface {
	FindByID(id int64) (*model.WorkDate, error)
}

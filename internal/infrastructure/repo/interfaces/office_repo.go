package interfaces

import "telegram-api/internal/domain/model"

type OfficeRepository interface {
	FindByID(id int64) (*model.Office, error)
	GetAll() ([]*model.Office, error)
}

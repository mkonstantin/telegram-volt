package interfaces

import "telegram-api/internal/domain/model"

type OfficeRepository interface {
	Create(office model.Office) (model.Office, error)
	Read(id int64) (model.Office, error)
	Update(office model.Office) (model.Office, error)
	Delete(id int64) error
}

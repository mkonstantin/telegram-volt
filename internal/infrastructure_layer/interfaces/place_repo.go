package interfaces

import "telegram-api/internal/domain_layer/model"

type PlaceRepository interface {
	Create(place model.Place) (model.Place, error)
	Read(id int64) (model.Place, error)
	Update(place model.Place) (model.Place, error)
	Delete(id int64) error
}

package repo

import "telegram-api/internal/model"

type PlaceRepository interface {
	Create(place model.Place) (model.Place, error)
	Read(id int64) (model.Place, error)
	Update(place model.Place) (model.Place, error)
	Delete(id int64) error
}

type PlaceRepositoryImpl struct {
	PlaceRepository PlaceRepository
}

func (s *PlaceRepositoryImpl) Create(place model.Place) (model.Place, error) {
	return model.Place{}, nil
}

func (s *PlaceRepositoryImpl) Read(id int64) (model.Place, error) {
	return model.Place{}, nil
}

func (s *PlaceRepositoryImpl) Update(place model.Place) (model.Place, error) {
	return model.Place{}, nil
}

func (s *PlaceRepositoryImpl) Delete(id int64) error {
	return nil
}

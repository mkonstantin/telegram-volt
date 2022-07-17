package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain_layer/model"
)

type placeRepositoryImpl struct {
	db *sqlx.DB
}

func NewPlaceRepository(db *sqlx.DB) model.PlaceRepository {
	return &placeRepositoryImpl{
		db: db,
	}
}

func (s *placeRepositoryImpl) Create(place model.Place) (model.Place, error) {
	return model.Place{}, nil
}

func (s *placeRepositoryImpl) Read(id int64) (model.Place, error) {
	return model.Place{}, nil
}

func (s *placeRepositoryImpl) Update(place model.Place) (model.Place, error) {
	return model.Place{}, nil
}

func (s *placeRepositoryImpl) Delete(id int64) error {
	return nil
}

package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/interface"
	repository "telegram-api/pkg"
)

type placeRepositoryImpl struct {
	db *sqlx.DB
}

func NewPlaceRepository(conn repository.Connection) _interface.PlaceRepository {
	return &placeRepositoryImpl{
		db: conn.Main,
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

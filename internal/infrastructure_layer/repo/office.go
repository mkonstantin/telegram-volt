package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain_layer/model"
)

type officeRepositoryImpl struct {
	db *sqlx.DB
}

func NewOfficeRepository(db *sqlx.DB) model.OfficeRepository {
	return &officeRepositoryImpl{
		db: db,
	}
}

func (s *officeRepositoryImpl) Create(office model.Office) (model.Office, error) {
	return model.Office{}, nil
}

func (s *officeRepositoryImpl) Read(id int64) (model.Office, error) {
	return model.Office{}, nil
}

func (s *officeRepositoryImpl) Update(office model.Office) (model.Office, error) {
	return model.Office{}, nil
}

func (s *officeRepositoryImpl) Delete(id int64) error {
	return nil
}

package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interfaces"
	repository "telegram-api/pkg"
)

type officeRepositoryImpl struct {
	db *sqlx.DB
}

func NewOfficeRepository(conn repository.Connection) interfaces.OfficeRepository {
	return &officeRepositoryImpl{
		db: conn.Main,
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

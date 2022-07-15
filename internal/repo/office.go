package repo

import "telegram-api/internal/model"

type OfficeRepository interface {
	Create(office model.Office) (model.Office, error)
	Read(id int64) (model.Office, error)
	Update(office model.Office) (model.Office, error)
	Delete(id int64) error
}

type OfficeRepositoryImpl struct {
	OfficeRepository OfficeRepository
}

func (s *OfficeRepositoryImpl) Create(office model.Office) (model.Office, error) {
	return model.Office{}, nil
}

func (s *OfficeRepositoryImpl) Read(id int64) (model.Office, error) {
	return model.Office{}, nil
}

func (s *OfficeRepositoryImpl) Update(office model.Office) (model.Office, error) {
	return model.Office{}, nil
}

func (s *OfficeRepositoryImpl) Delete(id int64) error {
	return nil
}

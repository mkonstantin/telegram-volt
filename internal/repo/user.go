package repo

import "telegram-api/internal/model"

type UserRepository interface {
	Create(user model.User) (model.User, error)
	Read(id int64) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id int64) error
}

type UserRepositoryImpl struct {
	UserRepository UserRepository
}

func (s *UserRepositoryImpl) Create(user model.User) (model.User, error) {
	return model.User{}, nil
}

func (s *UserRepositoryImpl) Read(id int64) (model.User, error) {
	return model.User{}, nil
}

func (s *UserRepositoryImpl) Update(user model.User) (model.User, error) {
	return model.User{}, nil
}

func (s *UserRepositoryImpl) Delete(id int64) error {
	return nil
}

package repo

import "telegram-api/internal/model"

type UserRepositoryImpl struct {
	UserRepository model.UserRepository
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

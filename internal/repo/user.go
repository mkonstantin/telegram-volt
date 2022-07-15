package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/model"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) model.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (s *userRepositoryImpl) Create(user model.User) (model.User, error) {
	return model.User{}, nil
}

func (s *userRepositoryImpl) Read(id int64) (model.User, error) {
	return model.User{}, nil
}

func (s *userRepositoryImpl) Update(user model.User) (model.User, error) {
	return model.User{}, nil
}

func (s *userRepositoryImpl) Delete(id int64) error {
	return nil
}

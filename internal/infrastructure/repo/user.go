package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/interface"
	repository "telegram-api/pkg"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(conn repository.Connection) _interface.UserRepository {
	return &userRepositoryImpl{
		db: conn.Main,
	}
}

func (s *userRepositoryImpl) GetByTelegramID(id int64) (model.User, error) {

	return model.User{}, nil
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

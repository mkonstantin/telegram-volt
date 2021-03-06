package repo

import (
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain_layer/model"
	"telegram-api/internal/infrastructure_layer/interfaces"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) interfaces.UserRepository {
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

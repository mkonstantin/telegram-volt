package repo

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
	repository "telegram-api/pkg"
)

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(conn repository.Connection) interfaces.UserRepository {
	return &userRepositoryImpl{
		db: conn.Main,
	}
}

func (s *userRepositoryImpl) GetByTelegramID(id int64) (*model.User, error) {
	sqQuery := sq.Select("*").
		From("user").
		Where(sq.Eq{"telegram_id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoU dto.User
	if err = s.db.Get(&dtoU, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dtoU.ToModel(), nil
}

func (s *userRepositoryImpl) Create(user model.User) error {
	sqQuery := sq.
		Insert("user").Columns("name", "telegram_id", "telegram_name").
		Values(user.Name, user.TelegramID, user.TelegramName)
	query, args, err := sqQuery.ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
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

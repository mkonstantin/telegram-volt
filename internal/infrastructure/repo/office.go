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

type officeRepositoryImpl struct {
	db *sqlx.DB
}

func NewOfficeRepository(conn repository.Connection) interfaces.OfficeRepository {
	return &officeRepositoryImpl{
		db: conn.Main,
	}
}

func (s *officeRepositoryImpl) FindByID(id int64) (*model.Office, error) {
	sqQuery := sq.Select("*").
		From("office").
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.Office
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *officeRepositoryImpl) GetAll() ([]*model.Office, error) {
	return []*model.Office{}, nil
}

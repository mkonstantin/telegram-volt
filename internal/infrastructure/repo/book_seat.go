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

type bookSeatRepositoryImpl struct {
	db *sqlx.DB
}

func NewBookSeatRepository(conn repository.Connection) interfaces.BookSeatRepository {
	return &bookSeatRepositoryImpl{
		db: conn.Main,
	}
}

func (s *bookSeatRepositoryImpl) FindByID(id int64) (*model.BookSeat, error) {
	sqQuery := sq.Select("*").
		From("seat").
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.BookSeat
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *bookSeatRepositoryImpl) GetAllByOfficeID(id int64) ([]*model.BookSeat, error) {
	sqQuery := sq.Select("*").
		From("seat").
		Where(sq.Eq{"office_id": id}).
		OrderBy("seat_number asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.BookSeat
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dto.ToBookSeatModels(dtoO), nil
}

func (s *bookSeatRepositoryImpl) BookSeatWithID(id int64) (*model.BookSeat, error) {
	//TODO implement me
	panic("implement me")
}

func (s *bookSeatRepositoryImpl) CancelBookSeatWithID(id int64) (*model.BookSeat, error) {
	//TODO implement me
	panic("implement me")
}

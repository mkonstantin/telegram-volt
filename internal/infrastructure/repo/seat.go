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

type seatRepositoryImpl struct {
	db *sqlx.DB
}

func NewSeatRepository(conn repository.Connection) interfaces.SeatRepository {
	return &seatRepositoryImpl{
		db: conn.Main,
	}
}

func (s *seatRepositoryImpl) FindByID(id int64) (*model.Seat, error) {
	sqQuery := sq.Select("*").
		From("seat").
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.Seat
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *seatRepositoryImpl) GetAllByOfficeID(id int64) ([]*model.Seat, error) {
	sqQuery := sq.Select("*").
		From("seat").
		Where(sq.Eq{"office_id": id}).
		OrderBy("seat_number asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.Seat
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dto.ToSeatModels(dtoO), nil
}

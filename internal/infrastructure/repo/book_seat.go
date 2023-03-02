package repo

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/common"
	"telegram-api/internal/infrastructure/repo/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
	repository "telegram-api/pkg"
	"time"
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
	sqQuery := sq.Select("bs.*, s1.seat_number, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.telegram_name, o1.name as office_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		InnerJoin("office as o1 ON bs.office_id=o1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.Eq{"bs.id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.BookSeat
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrBookSeatsNotFound
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *bookSeatRepositoryImpl) GetFreeSeatsByOfficeIDAndDate(id int64, dateStr string) ([]*model.BookSeat, error) {
	sqQuery := sq.Select("bs.*, s1.seat_number, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.telegram_name, o1.name as office_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		InnerJoin("office as o1 ON bs.office_id=o1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.And{sq.Eq{"bs.office_id": id}, sq.Eq{"bs.book_date": dateStr}, sq.Eq{"bs.user_id": nil}}).
		OrderBy("s1.seat_number asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.BookSeat
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrBookSeatsNotFound
		}
		return nil, err
	}

	return dto.ToBookSeatModels(dtoO), nil
}

func (s *bookSeatRepositoryImpl) GetAllByOfficeIDAndDate(id int64, dateStr string) ([]*model.BookSeat, error) {
	sqQuery := sq.Select("bs.*, s1.seat_number, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.telegram_name, o1.name as office_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		InnerJoin("office as o1 ON bs.office_id=o1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.And{sq.Eq{"bs.office_id": id}, sq.Eq{"bs.book_date": dateStr}}).
		OrderBy("s1.seat_number asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.BookSeat
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrBookSeatsNotFound
		}
		return nil, err
	}

	return dto.ToBookSeatModels(dtoO), nil
}

func (s *bookSeatRepositoryImpl) BookSeatWithID(userID, id int64) error {
	sqQuery := sq.Update("book_seat").
		Set("user_id", userID).
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *bookSeatRepositoryImpl) CancelBookSeatWithID(id int64) error {
	sqQuery := sq.Update("book_seat").
		Set("user_id", nil).
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *bookSeatRepositoryImpl) FindByUserID(userID int64) (*model.BookSeat, error) {
	sqQuery := sq.Select("bs.*, s1.seat_number, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.telegram_name, o1.name as office_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		InnerJoin("office as o1 ON bs.office_id=o1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.Eq{"bs.user_id": userID})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.BookSeat
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *bookSeatRepositoryImpl) FindByUserIDAndDate(userID int64, dateStr string) (*model.BookSeat, error) {
	sqQuery := sq.Select("bs.*, s1.seat_number, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.telegram_name, o1.name as office_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		InnerJoin("office as o1 ON bs.office_id=o1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.And{sq.Eq{"bs.user_id": userID}, sq.Eq{"bs.book_date": dateStr}})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.BookSeat
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *bookSeatRepositoryImpl) InsertSeat(officeID, seatID int64, dayDate time.Time) error {
	sqQuery := sq.
		Insert("book_seat").Columns("office_id", "seat_id", "book_date").
		Values(officeID, seatID, dayDate)
	query, args, err := sqQuery.ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

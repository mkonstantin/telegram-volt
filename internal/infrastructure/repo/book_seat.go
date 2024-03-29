package repo

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/common"
	"telegram-api/internal/infrastructure/repo/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
	repository "telegram-api/pkg"
	"telegram-api/pkg/tracing"
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

func (s *bookSeatRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name, o1.name as office_name, o1.time_zone").
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

func (s *bookSeatRepositoryImpl) GetFreeSeatsByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name, o1.name as office_name, o1.time_zone").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		InnerJoin("office as o1 ON bs.office_id=o1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.And{sq.Eq{"bs.office_id": id}, sq.Eq{"bs.book_date": dateStr}, sq.Eq{"bs.user_id": nil}, sq.Eq{"bs.hold": false}}).
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

func (s *bookSeatRepositoryImpl) FindByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.And{sq.Eq{"bs.office_id": id}, sq.Eq{"bs.book_date": dateStr}, sq.NotEq{"bs.user_id": nil}})

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

func (s *bookSeatRepositoryImpl) FindNotConfirmedByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name").
		From("book_seat as bs").
		InnerJoin("seat as s1 ON bs.seat_id = s1.id").
		LeftJoin("user as u1 ON bs.user_id = u1.id").
		Where(sq.And{sq.Eq{"bs.office_id": id}, sq.Eq{"bs.book_date": dateStr}, sq.NotEq{"bs.user_id": nil}, sq.Eq{"bs.confirm": false}})

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

func (s *bookSeatRepositoryImpl) GetAllByOfficeIDAndDate(ctx context.Context, id int64, dateStr string) ([]*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name, o1.name as office_name, o1.time_zone").
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

func (s *bookSeatRepositoryImpl) BookSeatWithID(ctx context.Context, userID, id int64, confirm bool) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("book_seat").
		Set("user_id", userID).
		Set("confirm", confirm).
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *bookSeatRepositoryImpl) CancelBookSeatWithID(ctx context.Context, id int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

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

func (s *bookSeatRepositoryImpl) FindByUserID(ctx context.Context, userID int64) (*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name, o1.name as office_name, o1.time_zone").
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

func (s *bookSeatRepositoryImpl) FindByUserIDAndDate(ctx context.Context, userID int64, dateStr string) (*model.BookSeat, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("bs.*, s1.seat_sign, s1.have_monitor, u1.name as user_name, u1.telegram_id, u1.chat_id, u1.telegram_name, o1.name as office_name, o1.time_zone").
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

func (s *bookSeatRepositoryImpl) InsertSeat(ctx context.Context, officeID, seatID int64, dayDate time.Time) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

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

func (s *bookSeatRepositoryImpl) ConfirmBookSeat(ctx context.Context, seatID int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("book_seat").
		Set("confirm", true).
		Where(sq.Eq{"id": seatID})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *bookSeatRepositoryImpl) HoldSeatWithID(ctx context.Context, id int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("book_seat").
		Set("hold", true).
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *bookSeatRepositoryImpl) CancelHoldSeatWithID(ctx context.Context, id int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("book_seat").
		Set("hold", false).
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

package repo

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
	repository "telegram-api/pkg"
	"telegram-api/pkg/tracing"
	"time"
)

type workDateRepositoryImpl struct {
	db *sqlx.DB
}

func NewWorkDateRepository(conn repository.Connection) interfaces.WorkDateRepository {
	return &workDateRepositoryImpl{
		db: conn.Main,
	}
}

func (s *workDateRepositoryImpl) DoneAllPastByDate(ctx context.Context, date string) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("work_date as wd").
		Set("status", model.StatusDone).
		Where(sq.And{sq.NotEq{"wd.status": model.StatusDone}, sq.Lt{"wd.work_date": date}})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *workDateRepositoryImpl) FindByDates(ctx context.Context, startDate string, endDate string) ([]model.WorkDate, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("work_date as wd").
		Where(sq.And{sq.GtOrEq{"wd.work_date": startDate}, sq.Lt{"wd.work_date": endDate}}).
		OrderBy("wd.work_date asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.WorkDate
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(dtoO) == 0 {
		return nil, nil
	}
	return dto.ToWorkDateModels(dtoO), nil
}

func (s *workDateRepositoryImpl) FindByDatesAndStatus(ctx context.Context, startDate string, endDate string,
	status model.DateStatus) ([]model.WorkDate, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("work_date as wd").
		Where(sq.And{sq.Eq{"wd.status": status},
			sq.GtOrEq{"wd.work_date": startDate},
			sq.Lt{"wd.work_date": endDate}}).
		OrderBy("wd.work_date asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.WorkDate
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(dtoO) == 0 {
		return nil, nil
	}
	return dto.ToWorkDateModels(dtoO), nil
}

func (s *workDateRepositoryImpl) FindByStatus(ctx context.Context, status model.DateStatus) ([]model.WorkDate, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("work_date as wd").
		Where(sq.Eq{"wd.status": status}).
		OrderBy("wd.work_date asc")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.WorkDate
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(dtoO) == 0 {
		return nil, nil
	}
	return dto.ToWorkDateModels(dtoO), nil
}

func (s *workDateRepositoryImpl) GetLastByDate(ctx context.Context) (*model.WorkDate, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("work_date as wd").
		OrderBy("wd.work_date desc").
		Limit(1)

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.WorkDate
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(dtoO) == 0 {
		return nil, nil
	}
	return dtoO[0].ToModel(), nil
}

func (s *workDateRepositoryImpl) InsertDate(ctx context.Context, date time.Time) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.
		Insert("work_date").Columns("status", "work_date").
		Values(model.StatusWait, date)
	query, args, err := sqQuery.ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *workDateRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.WorkDate, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("work_date").
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO dto.WorkDate
	if err = s.db.Get(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *workDateRepositoryImpl) UpdateStatusByID(ctx context.Context, id int64, status model.DateStatus) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("work_date").
		Set("status", status).
		Where(sq.Eq{"id": id})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

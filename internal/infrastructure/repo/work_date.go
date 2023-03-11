package repo

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/repo/dto"
	"telegram-api/internal/infrastructure/repo/interfaces"
	repository "telegram-api/pkg"
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

func (s *workDateRepositoryImpl) GetLastByDate() (*model.WorkDate, error) {
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

func (s *workDateRepositoryImpl) InsertDate(date time.Time) error {
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

func (s *workDateRepositoryImpl) FindByID(id int64) (*model.WorkDate, error) {
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

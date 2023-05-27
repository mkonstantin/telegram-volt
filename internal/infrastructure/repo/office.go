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
)

type officeRepositoryImpl struct {
	db *sqlx.DB
}

func NewOfficeRepository(conn repository.Connection) interfaces.OfficeRepository {
	return &officeRepositoryImpl{
		db: conn.Main,
	}
}

func (s *officeRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Office, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

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
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dtoO.ToModel(), nil
}

func (s *officeRepositoryImpl) GetAll(ctx context.Context) ([]*model.Office, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("office")

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.Office
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dto.ToOfficeModels(dtoO), nil
}

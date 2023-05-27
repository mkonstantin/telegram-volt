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

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(conn repository.Connection) interfaces.UserRepository {
	return &userRepositoryImpl{
		db: conn.Main,
	}
}

func (s *userRepositoryImpl) GetUsersToNotify(ctx context.Context, notifyOfficeID int64) ([]*model.User, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Select("*").
		From("user").
		Where(sq.And{sq.Eq{"notify_office_id": notifyOfficeID}, sq.NotEq{"chat_id": 0}})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	var dtoO []dto.User
	if err = s.db.Select(&dtoO, query, args...); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			return nil, nil
		}
		return nil, err
	}

	return dto.ToUserModels(dtoO), nil
}

func (s *userRepositoryImpl) GetByTelegramID(ctx context.Context, id int64) (*model.User, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

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

func (s *userRepositoryImpl) Create(ctx context.Context, user model.User) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.
		Insert("user").Columns("name", "telegram_id", "telegram_name", "chat_id").
		Values(user.Name, user.TelegramID, user.TelegramName, user.ChatID)
	query, args, err := sqQuery.ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *userRepositoryImpl) SetChatID(ctx context.Context, chatID, tgID int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("user").
		Set("chat_id", chatID).
		Where(sq.Eq{"telegram_id": tgID})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *userRepositoryImpl) SetOffice(ctx context.Context, officeID, tgID int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("user").
		Set("office_id", officeID).
		Where(sq.Eq{"telegram_id": tgID})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *userRepositoryImpl) Subscribe(ctx context.Context, tgID, officeID int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("user").
		Set("notify_office_id", officeID).
		Where(sq.Eq{"telegram_id": tgID})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *userRepositoryImpl) Unsubscribe(ctx context.Context, tgID int64) error {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	sqQuery := sq.Update("user").
		Set("notify_office_id", 0).
		Where(sq.Eq{"telegram_id": tgID})

	query, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

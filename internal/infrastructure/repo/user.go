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

func (s *userRepositoryImpl) GetUsersToNotify(notifyOfficeID int64) ([]*model.User, error) {
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
		Insert("user").Columns("name", "telegram_id", "telegram_name", "chat_id").
		Values(user.Name, user.TelegramID, user.TelegramName, user.ChatID)
	query, args, err := sqQuery.ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(query, args...)
	return nil
}

func (s *userRepositoryImpl) SetChatID(chatID, tgID int64) error {
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

func (s *userRepositoryImpl) SetOffice(officeID, tgID int64) error {
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

func (s *userRepositoryImpl) Subscribe(tgID, officeID int64) error {
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

func (s *userRepositoryImpl) Unsubscribe(tgID int64) error {
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

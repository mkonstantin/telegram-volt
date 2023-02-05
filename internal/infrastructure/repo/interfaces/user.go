package interfaces

import "telegram-api/internal/domain/model"

type UserRepository interface {
	GetByTelegramID(id int64) (*model.User, error)
	Create(user model.User) error
	SetOffice(officeID, tgID int64) error
	Subscribe(tgID, officeID int64) error
	Unsubscribe(tgID int64) error
}

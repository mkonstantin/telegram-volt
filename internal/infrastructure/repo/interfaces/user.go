package interfaces

import "telegram-api/internal/domain/model"

type UserRepository interface {
	GetUsersToNotify(notifyOfficeID int64) ([]*model.User, error)
	GetByTelegramID(id int64) (*model.User, error)
	Create(user model.User) error
	SetChatID(chatID, tgID int64) error
	SetOffice(officeID, tgID int64) error
	Subscribe(tgID, officeID int64) error
	Unsubscribe(tgID int64) error
}

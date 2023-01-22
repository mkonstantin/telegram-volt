package interfaces

import "telegram-api/internal/domain/model"

type UserRepository interface {
	GetByTelegramID(id int64) (*model.User, error)
	Create(user model.User) error
}

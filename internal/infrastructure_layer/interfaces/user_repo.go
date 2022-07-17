package interfaces

import "telegram-api/internal/domain_layer/model"

type UserRepository interface {
	Create(user model.User) (model.User, error)
	Read(id int64) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id int64) error
}

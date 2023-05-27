package interfaces

import (
	"context"
	"telegram-api/internal/domain/model"
)

type UserRepository interface {
	GetUsersToNotify(ctx context.Context, notifyOfficeID int64) ([]*model.User, error)
	GetByTelegramID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user model.User) error
	SetChatID(ctx context.Context, chatID, tgID int64) error
	SetOffice(ctx context.Context, officeID, tgID int64) error
	Subscribe(ctx context.Context, tgID, officeID int64) error
	Unsubscribe(ctx context.Context, tgID int64) error
}

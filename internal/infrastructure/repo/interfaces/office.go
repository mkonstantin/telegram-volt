package interfaces

import (
	"context"
	"telegram-api/internal/domain/model"
)

type OfficeRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Office, error)
	GetAll(ctx context.Context) ([]*model.Office, error)
}

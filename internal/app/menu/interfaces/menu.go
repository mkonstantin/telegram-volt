package interfaces

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OfficeMenu interface {
	Call(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

type OfficeListMenu interface {
	Call(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

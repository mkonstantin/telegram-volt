package interfaces

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type OfficeMenu interface {
	Call(ctx context.Context, title string, officeID int64) (*tgbotapi.MessageConfig, error)
}

type OfficeListMenu interface {
	Call(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

type SeatListMenu interface {
	Call(ctx context.Context, date time.Time, officeID int64) (*tgbotapi.MessageConfig, error)
}

type FreeSeatMenu interface {
	Call(ctx context.Context, bookSeatID int64) (*tgbotapi.MessageConfig, error)
}

type OwnSeatMenu interface {
	Call(ctx context.Context, bookSeatID int64) (*tgbotapi.MessageConfig, error)
}

type DateMenu interface {
	Call(ctx context.Context) (*tgbotapi.MessageConfig, error)
}

type HoldSeatMenu interface {
	Call(ctx context.Context, bookSeatID int64) (*tgbotapi.MessageConfig, error)
}

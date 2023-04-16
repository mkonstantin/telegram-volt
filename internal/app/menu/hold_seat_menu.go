package menu

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/form"
	"telegram-api/internal/app/menu/interfaces"
)

type holdSeatMenuImpl struct {
	holdSeatForm form.HoldSeatForm
	logger       *zap.Logger
}

func NewHoldSeatMenu(
	holdSeatForm form.HoldSeatForm,
	logger *zap.Logger) interfaces.HoldSeatMenu {

	return &holdSeatMenuImpl{
		holdSeatForm: holdSeatForm,
		logger:       logger,
	}
}

func (f *holdSeatMenuImpl) Call(ctx context.Context, bookSeatID int64) (*tgbotapi.MessageConfig, error) {
	message := "Место закреплено, хотите снять закрепление?"

	formData := form.HoldSeatFormData{
		Message:    message,
		BookSeatID: bookSeatID,
	}
	return f.holdSeatForm.Build(ctx, formData)
}

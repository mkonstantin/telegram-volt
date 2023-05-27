package form

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/router/constants"
	"telegram-api/pkg/tracing"
)

type HoldSeatFormData struct {
	Message    string
	BookSeatID int64
}

type HoldSeatForm interface {
	Build(ctx context.Context, data HoldSeatFormData) (*tgbotapi.MessageConfig, error)
}

type holdSeatFormImpl struct {
	logger *zap.Logger
}

func NewHoldSeatForm(logger *zap.Logger) HoldSeatForm {
	return &holdSeatFormImpl{
		logger: logger,
	}
}

func (f *holdSeatFormImpl) Build(ctx context.Context, data HoldSeatFormData) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = data.Message

	b1 := &dto.InlineRequest{
		Type:       constants.HoldSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionCancelHoldYes,
	}
	b2 := &dto.InlineRequest{
		Type:       constants.HoldSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionCancelHoldNo,
	}

	butt1, err := json.Marshal(b1)
	if err != nil {
		return nil, err
	}
	butt2, err := json.Marshal(b2)
	if err != nil {
		return nil, err
	}

	button1 := tgbotapi.NewInlineKeyboardButtonData("Снять", string(butt1))
	button2 := tgbotapi.NewInlineKeyboardButtonData("К списку мест", string(butt2))

	row := tgbotapi.NewInlineKeyboardRow(button1, button2)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	msg.ReplyMarkup = keyboard

	return &msg, nil
}

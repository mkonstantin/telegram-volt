package form

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/router/constants"
)

type FreeSeatFormData struct {
	Message    string
	BookSeatID int64
}

type FreeSeatForm interface {
	Build(ctx context.Context, data FreeSeatFormData) (*tgbotapi.MessageConfig, error)
}

type freeSeatFormImpl struct {
	logger *zap.Logger
}

func NewFreeSeatForm(logger *zap.Logger) FreeSeatForm {
	return &freeSeatFormImpl{
		logger: logger,
	}
}

func (f *freeSeatFormImpl) Build(ctx context.Context, data FreeSeatFormData) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = data.Message

	b1 := &dto.InlineRequest{
		Type:       constants.FreeSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionBookYes,
	}
	b2 := &dto.InlineRequest{
		Type:       constants.FreeSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionBookNo,
	}

	butt1, err := json.Marshal(b1)
	if err != nil {
		return nil, err
	}
	butt2, err := json.Marshal(b2)
	if err != nil {
		return nil, err
	}

	button1 := tgbotapi.NewInlineKeyboardButtonData("Занять", string(butt1))
	button2 := tgbotapi.NewInlineKeyboardButtonData("К списку мест", string(butt2))
	row := tgbotapi.NewInlineKeyboardRow(button1, button2)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	msg.ReplyMarkup = keyboard

	return &msg, nil
}

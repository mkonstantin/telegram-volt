package form

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/config"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/router/constants"
)

type OwnSeatFormData struct {
	Message    string
	BookSeatID int64
}

type OwnSeatForm interface {
	Build(ctx context.Context, data OwnSeatFormData) (*tgbotapi.MessageConfig, error)
}

type ownSeatFormImpl struct {
	cfg    config.AppConfig
	logger *zap.Logger
}

func NewOwnSeatForm(
	cfg config.AppConfig,
	logger *zap.Logger) OwnSeatForm {
	return &ownSeatFormImpl{
		cfg:    cfg,
		logger: logger,
	}
}

func (f *ownSeatFormImpl) Build(ctx context.Context, data OwnSeatFormData) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = data.Message

	b1 := &dto.InlineRequest{
		Type:       constants.OwnSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionCancelBookYes,
	}
	b2 := &dto.InlineRequest{
		Type:       constants.OwnSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionCancelBookNo,
	}

	butt1, err := json.Marshal(b1)
	if err != nil {
		return nil, err
	}
	butt2, err := json.Marshal(b2)
	if err != nil {
		return nil, err
	}

	button1 := tgbotapi.NewInlineKeyboardButtonData("Освободить", string(butt1))
	button2 := tgbotapi.NewInlineKeyboardButtonData("К списку мест", string(butt2))

	currentUser := model.GetCurrentUser(ctx)

	row := tgbotapi.NewInlineKeyboardRow(button1, button2)

	if f.cfg.IsAdmin(currentUser.TelegramName) {
		buttonAdmin, err := f.addAdminStuff(data)
		if err != nil {
			return nil, err
		}
		if buttonAdmin != nil {
			row = append(row, *buttonAdmin)
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	msg.ReplyMarkup = keyboard

	return &msg, nil
}

func (f *ownSeatFormImpl) addAdminStuff(data OwnSeatFormData) (*tgbotapi.InlineKeyboardButton, error) {
	b := &dto.InlineRequest{
		Type:       constants.OwnSeatMenuTap,
		BookSeatID: data.BookSeatID,
		Action:     dto.ActionCancelHold,
	}

	butt, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	button := tgbotapi.NewInlineKeyboardButtonData("Снять закрепление", string(butt))
	return &button, nil
}

package form

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/router/constants"
	"time"
)

type InfoFormData struct {
	ChatID  int64
	Message string
	Office  *model.Office
	Date    time.Time
}

type InfoMenuForm interface {
	Build(ctx context.Context, data InfoFormData) (*tgbotapi.MessageConfig, error)
}

type infoMenuFormImpl struct {
	logger *zap.Logger
}

func NewInfoMenuForm(logger *zap.Logger) InfoMenuForm {
	return &infoMenuFormImpl{
		logger: logger,
	}
}

func (f *infoMenuFormImpl) Build(_ context.Context, data InfoFormData) (*tgbotapi.MessageConfig, error) {

	msg := tgbotapi.NewMessage(data.ChatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton

	resp := &dto.InlineRequest{
		Type:     constants.InformerTap,
		BookDate: &data.Date,
	}
	responseData, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	button := tgbotapi.NewInlineKeyboardButtonData("Перейти", string(responseData))
	row := tgbotapi.NewInlineKeyboardRow(button)
	rows = append(rows, row)

	var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)

	msg.Text = data.Message
	msg.ReplyMarkup = chooseOfficeKeyboard

	return &msg, nil
}

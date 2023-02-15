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

type OfficeListFormData struct {
	Message string
	Offices []*model.Office
}

type OfficeListForm interface {
	Build(ctx context.Context, data OfficeListFormData) (*tgbotapi.MessageConfig, error)
}

type officeListMenuFormImpl struct {
	logger *zap.Logger
}

func NewOfficeListForm(logger *zap.Logger) OfficeListForm {
	return &officeListMenuFormImpl{
		logger: logger,
	}
}

func (o officeListMenuFormImpl) Build(ctx context.Context, data OfficeListFormData) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, office := range data.Offices {
		resp := &dto.InlineRequest{
			Type:     constants.OfficeListTap,
			OfficeID: office.ID,
		}
		responseData, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}

		button := tgbotapi.NewInlineKeyboardButtonData(office.Name, string(responseData))
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}

	var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)

	msg.Text = data.Message
	msg.ReplyMarkup = chooseOfficeKeyboard

	return &msg, nil
}

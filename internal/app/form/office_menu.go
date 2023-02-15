package form

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler/dto"
	"telegram-api/internal/infrastructure/router"
)

type OfficeMenuForm interface {
	Build(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
}

type officeMenuFormImpl struct {
	logger *zap.Logger
}

func NewOfficeMenuForm(logger *zap.Logger) OfficeMenuForm {
	return &officeMenuFormImpl{
		logger: logger,
	}
}

func (o *officeMenuFormImpl) Build(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	b1 := &dto.InlineRequest{
		Type:     router.OfficeMenuTap,
		OfficeID: result.Office.ID,
		Action:   dto.OfficeMenuFreeSeats,
	}
	b2 := &dto.InlineRequest{
		Type:     router.OfficeMenuTap,
		OfficeID: result.Office.ID,
		Action:   dto.OfficeMenuSubscribe,
	}
	b3 := &dto.InlineRequest{
		Type:     router.OfficeMenuTap,
		OfficeID: result.Office.ID,
		Action:   dto.OfficeMenuChooseAnotherOffice,
	}

	butt1, err := json.Marshal(b1)
	if err != nil {
		return nil, err
	}
	butt2, err := json.Marshal(b2)
	if err != nil {
		return nil, err
	}
	butt3, err := json.Marshal(b3)
	if err != nil {
		return nil, err
	}

	button1 := tgbotapi.NewInlineKeyboardButtonData("Показать места", string(butt1))
	button2 := tgbotapi.NewInlineKeyboardButtonData(result.SubscribeButtonText, string(butt2))
	button3 := tgbotapi.NewInlineKeyboardButtonData("Выбрать другой офис", string(butt3))
	row1 := tgbotapi.NewInlineKeyboardRow(button1)
	row2 := tgbotapi.NewInlineKeyboardRow(button2)
	row3 := tgbotapi.NewInlineKeyboardRow(button3)
	confirmOfficeKeyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)

	msg.Text = result.Message
	msg.ReplyMarkup = confirmOfficeKeyboard
	//msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

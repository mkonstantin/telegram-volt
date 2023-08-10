package form

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/helper"
	"telegram-api/internal/infrastructure/router/constants"
	"telegram-api/pkg/tracing"
)

type OfficeMenuFormData struct {
	Message             string
	SubscribeButtonText string
	Office              *model.Office
	BookSeats           []*model.BookSeat
	NeedConfirmBookSeat *model.BookSeat
}

type OfficeMenuForm interface {
	Build(ctx context.Context, data OfficeMenuFormData) (*tgbotapi.MessageConfig, error)
}

type officeMenuFormImpl struct {
	logger *zap.Logger
}

func NewOfficeMenuForm(logger *zap.Logger) OfficeMenuForm {
	return &officeMenuFormImpl{
		logger: logger,
	}
}

func (o *officeMenuFormImpl) Build(ctx context.Context, data OfficeMenuFormData) (*tgbotapi.MessageConfig, error) {
	ctx, span, _ := tracing.StartSpan(ctx, tracing.GetSpanName())
	defer span.End()

	var sum [][]tgbotapi.InlineKeyboardButton

	if data.NeedConfirmBookSeat != nil {
		b := &dto.InlineRequest{
			Type:       constants.OfficeMenuTap,
			Action:     dto.OfficeMenuConfirm,
			BookSeatID: data.NeedConfirmBookSeat.ID,
		}
		butt, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}
		buttonMessage := fmt.Sprintf("ПОДТВЕРДИТЬ бронь на сегодня")
		button := tgbotapi.NewInlineKeyboardButtonData(buttonMessage, string(butt))
		row := tgbotapi.NewInlineKeyboardRow(button)
		sum = append(sum, row)
	}

	sum, err := addStandartButtons(sum, data)
	if err != nil {
		return nil, err
	}

	for _, seat := range data.BookSeats {
		b := &dto.InlineRequest{
			Type:       constants.OfficeMenuTap,
			Action:     dto.OfficeMenuCancelBook,
			BookSeatID: seat.ID,
		}
		butt, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}
		buttonMessage := fmt.Sprintf("ОТМЕНИТЬ бронь на %s", seat.BookDate.Format(helper.DateFormat))
		button := tgbotapi.NewInlineKeyboardButtonData(buttonMessage, string(butt))
		row := tgbotapi.NewInlineKeyboardRow(button)
		sum = append(sum, row)
	}

	chatID := model.GetCurrentChatID(ctx)
	msg := tgbotapi.NewMessage(chatID, "")
	confirmOfficeKeyboard := tgbotapi.NewInlineKeyboardMarkup(sum...)
	msg.Text = data.Message
	msg.ReplyMarkup = confirmOfficeKeyboard

	return &msg, nil
}

func addStandartButtons(sum [][]tgbotapi.InlineKeyboardButton, data OfficeMenuFormData) ([][]tgbotapi.InlineKeyboardButton, error) {

	b1 := &dto.InlineRequest{
		Type:     constants.OfficeMenuTap,
		OfficeID: data.Office.ID,
		Action:   dto.OfficeMenuFreeSeats,
	}
	b2 := &dto.InlineRequest{
		Type:     constants.OfficeMenuTap,
		OfficeID: data.Office.ID,
		Action:   dto.OfficeMenuSubscribe,
	}
	b3 := &dto.InlineRequest{
		Type:     constants.OfficeMenuTap,
		OfficeID: data.Office.ID,
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
	button2 := tgbotapi.NewInlineKeyboardButtonData(data.SubscribeButtonText, string(butt2))
	button3 := tgbotapi.NewInlineKeyboardButtonData("⬅️ Выбрать другой хотдеск", string(butt3))
	row1 := tgbotapi.NewInlineKeyboardRow(button1)
	row2 := tgbotapi.NewInlineKeyboardRow(button2)
	row3 := tgbotapi.NewInlineKeyboardRow(button3)

	sum = append(sum, row1)
	sum = append(sum, row2)
	sum = append(sum, row3)

	return sum, nil
}

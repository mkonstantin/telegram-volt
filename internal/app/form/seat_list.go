package form

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/handler/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/router/constants"
)

type SeatListFormData struct {
	Message   string
	BookSeats []*model.BookSeat
}

type SeatListForm interface {
	Build(ctx context.Context, data SeatListFormData) (*tgbotapi.MessageConfig, error)
}

type seatListFormImpl struct {
	logger *zap.Logger
}

func NewSeatListForm(logger *zap.Logger) SeatListForm {
	return &seatListFormImpl{
		logger: logger,
	}
}

func (o seatListFormImpl) Build(ctx context.Context, data SeatListFormData) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, bookSeat := range data.BookSeats {
		resp := &dto.InlineRequest{
			Type:       constants.SeatListTap,
			BookSeatID: bookSeat.ID,
		}
		responseData, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}

		var button tgbotapi.InlineKeyboardButton
		if bookSeat.User != nil {
			str := fmt.Sprintf("Место %s, занято: %s", bookSeat.Seat.SeatSign, bookSeat.User.Name)
			button = tgbotapi.NewInlineKeyboardButtonData(str, string(responseData))
		} else {
			if bookSeat.Hold {
				str := fmt.Sprintf("Место %s. Временно закреплено!", bookSeat.Seat.SeatSign)
				button = tgbotapi.NewInlineKeyboardButtonData(str, string(responseData))
			} else {
				str := fmt.Sprintf("Место %s. Свободно!", bookSeat.Seat.SeatSign)
				button = tgbotapi.NewInlineKeyboardButtonData(str, string(responseData))
			}
		}

		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}

	var msg tgbotapi.MessageConfig

	if len(data.BookSeats) > 0 {
		respBack := &dto.InlineRequest{
			Type:   constants.SeatListTap,
			Action: dto.Back,
		}
		backData, err := json.Marshal(respBack)
		if err != nil {
			return nil, err
		}
		button := tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", string(backData))
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)

		msg = tgbotapi.NewMessage(chatID, "")
		var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			rows...,
		)
		msg.Text = data.Message
		msg.ReplyMarkup = chooseOfficeKeyboard
	} else {
		msg = tgbotapi.NewMessage(chatID, "Мест не найдено")
	}

	return &msg, nil
}

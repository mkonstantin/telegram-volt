package form

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler/dto"
	"telegram-api/internal/infrastructure/router/constants"
)

type SeatListForm interface {
	Build(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
}

type seatListFormImpl struct {
	logger *zap.Logger
}

func NewSeatListForm(logger *zap.Logger) SeatListForm {
	return &seatListFormImpl{
		logger: logger,
	}
}

func (o seatListFormImpl) Build(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, bookSeat := range result.BookSeats {
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
			str := fmt.Sprintf("Место %d, занято: %s", bookSeat.Seat.SeatNumber, bookSeat.User.Name)
			button = tgbotapi.NewInlineKeyboardButtonData(str, string(responseData))
		} else {
			str := fmt.Sprintf("Место %d. Свободно!", bookSeat.Seat.SeatNumber)
			button = tgbotapi.NewInlineKeyboardButtonData(str, string(responseData))
		}

		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}

	var msg tgbotapi.MessageConfig

	if len(result.BookSeats) > 0 {
		msg = tgbotapi.NewMessage(chatID, "")
		var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			rows...,
		)
		msg.Text = result.Message
		msg.ReplyMarkup = chooseOfficeKeyboard
	} else {
		msg = tgbotapi.NewMessage(chatID, "Мест не найдено")
	}

	return &msg, nil
}

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
	"time"
)

type DaySeat struct {
	Date        time.Time
	SeatsNumber int
}

type DateMenuFormData struct {
	Message     string
	Office      *model.Office
	SeatByDates []DaySeat
}

type DateMenuForm interface {
	Build(ctx context.Context, data DateMenuFormData) (*tgbotapi.MessageConfig, error)
}

type freeDateMenuFormImpl struct {
	logger *zap.Logger
}

func NewDateMenutForm(logger *zap.Logger) DateMenuForm {
	return &freeDateMenuFormImpl{
		logger: logger,
	}
}

func (f *freeDateMenuFormImpl) Build(ctx context.Context, data DateMenuFormData) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton

	respMap := &dto.InlineRequest{
		Type:   constants.DateMenuTap,
		Action: dto.DateListShowMap,
	}
	mapData, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}
	buttonMap := tgbotapi.NewInlineKeyboardButtonData("📍 Карта расположения мест 📍", string(mapData))
	rowMap := tgbotapi.NewInlineKeyboardRow(buttonMap)
	rows = append(rows, rowMap)

	for _, seatByDate := range data.SeatByDates {
		resp := &dto.InlineRequest{
			Type:     constants.DateMenuTap,
			BookDate: &seatByDate.Date,
		}
		responseData, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}

		formattedDate := seatByDate.Date.Format(helper.DateFormat)

		text := fmt.Sprintf("%s : %d свободных мест", formattedDate, seatByDate.SeatsNumber)
		button := tgbotapi.NewInlineKeyboardButtonData(text, string(responseData))
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}

	respBack := &dto.InlineRequest{
		Type:   constants.DateMenuTap,
		Action: dto.Back,
	}
	backData, err := json.Marshal(respBack)
	if err != nil {
		return nil, err
	}
	button := tgbotapi.NewInlineKeyboardButtonData("Назад", string(backData))
	row := tgbotapi.NewInlineKeyboardRow(button)
	rows = append(rows, row)

	var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)

	msg.Text = data.Message
	msg.ReplyMarkup = chooseOfficeKeyboard

	return &msg, nil
}

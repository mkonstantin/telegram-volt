package handler

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/infrastructure/handler/dto"
)

type MessageFormer interface {
	FormChooseOfficeMenuMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormOfficeMenuMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
}

type messageFormerImpl struct {
	logger *zap.Logger
}

func NewMessageFormer(logger *zap.Logger) MessageFormer {
	return &messageFormerImpl{
		logger: logger,
	}
}

func (s *messageFormerImpl) Start(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {

	return tgbotapi.MessageConfig{}, nil
}

func (s *messageFormerImpl) FormChooseOfficeMenuMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, office := range result.Offices {
		resp := &dto.CommandResponse{
			Type:     usecase.ChooseOfficeMenu,
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

	msg.Text = result.Message
	msg.ReplyMarkup = chooseOfficeKeyboard
	//msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

func (s *messageFormerImpl) FormOfficeMenuMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")

	b1 := &dto.CommandResponse{
		Type:     usecase.OfficeMenu,
		OfficeID: result.Office.ID,
		Action:   dto.OfficeMenuFreeSeats,
	}
	b2 := &dto.CommandResponse{
		Type:     usecase.OfficeMenu,
		OfficeID: result.Office.ID,
		Action:   dto.OfficeMenuSubscribe,
	}
	b3 := &dto.CommandResponse{
		Type:     usecase.OfficeMenu,
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

	button1 := tgbotapi.NewInlineKeyboardButtonData("Свободные места", string(butt1))
	button2 := tgbotapi.NewInlineKeyboardButtonData("Подписаться на запись", string(butt2))
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

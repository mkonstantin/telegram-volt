package handler

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/infrastructure/handler/dto"
)

type MessageFormer interface {
	FormChooseOfficeMenuMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormOfficeMenuMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormSeatListMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormBookSeatMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
}

type messageFormerImpl struct {
	logger *zap.Logger
}

func NewMessageFormer(logger *zap.Logger) MessageFormer {
	return &messageFormerImpl{
		logger: logger,
	}
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

	button1 := tgbotapi.NewInlineKeyboardButtonData("Показать места", string(butt1))
	button2 := tgbotapi.NewInlineKeyboardButtonData("Подписаться на свободные места", string(butt2))
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

func (s *messageFormerImpl) FormSeatListMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, bookSeat := range result.BookSeats {
		resp := &dto.CommandResponse{
			Type:       usecase.ChooseSeatsMenu,
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
		msg = tgbotapi.NewMessage(result.ChatID, "")
		var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			rows...,
		)
		msg.Text = result.Message
		msg.ReplyMarkup = chooseOfficeKeyboard
	} else {
		msg = tgbotapi.NewMessage(result.ChatID, "Мест не найдено")
	}

	return &msg, nil
}

func (s *messageFormerImpl) FormBookSeatMsg(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")
	switch result.Key {
	case usecase.SeatOwn:
		msg.Text = result.Message

		b1 := &dto.CommandResponse{
			Type:       usecase.SeatOwn,
			BookSeatID: result.BookSeatID,
			Action:     dto.ActionCancelBookYes,
		}
		b2 := &dto.CommandResponse{
			Type:       usecase.SeatOwn,
			BookSeatID: result.BookSeatID,
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

		button1 := tgbotapi.NewInlineKeyboardButtonData("Да", string(butt1))
		button2 := tgbotapi.NewInlineKeyboardButtonData("Нет", string(butt2))
		row := tgbotapi.NewInlineKeyboardRow(button1, button2)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
		msg.ReplyMarkup = keyboard

	case usecase.SeatBusy:
		msg.Text = result.Message
	case usecase.SeatFree:
		msg.Text = result.Message
	}
	return &msg, nil
}

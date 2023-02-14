package handler

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler/dto"
	"telegram-api/internal/infrastructure/router"
)

type MessageFormer interface {
	FormChooseOfficeMenuMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormOfficeMenuMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormSeatListMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormBookSeatMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormBookSeatResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormCancelBookResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
	FormTextMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error)
}

type messageFormerImpl struct {
	logger *zap.Logger
}

func NewMessageFormer(logger *zap.Logger) MessageFormer {
	return &messageFormerImpl{
		logger: logger,
	}
}

func (s *messageFormerImpl) FormChooseOfficeMenuMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, office := range result.Offices {
		resp := &dto.InlineRequest{
			Type:     router.OfficeListTap,
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

func (s *messageFormerImpl) FormOfficeMenuMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

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

func (s *messageFormerImpl) FormSeatListMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, bookSeat := range result.BookSeats {
		resp := &dto.InlineRequest{
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

func (s *messageFormerImpl) FormBookSeatMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")
	switch result.Key {
	case usecase.SeatOwn:
		msg.Text = result.Message

		b1 := &dto.InlineRequest{
			Type:       usecase.SeatOwn,
			BookSeatID: result.BookSeatID,
			Action:     dto.ActionCancelBookYes,
		}
		b2 := &dto.InlineRequest{
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

		button1 := tgbotapi.NewInlineKeyboardButtonData("Освободить", string(butt1))
		button2 := tgbotapi.NewInlineKeyboardButtonData("К списку мест", string(butt2))
		row := tgbotapi.NewInlineKeyboardRow(button1, button2)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
		msg.ReplyMarkup = keyboard

	case usecase.SeatBusy:
		msg.Text = result.Message
	case usecase.SeatFree:
		msg.Text = result.Message

		b1 := &dto.InlineRequest{
			Type:       usecase.SeatFree,
			BookSeatID: result.BookSeatID,
			Action:     dto.ActionBookYes,
		}
		b2 := &dto.InlineRequest{
			Type:       usecase.SeatFree,
			BookSeatID: result.BookSeatID,
			Action:     dto.ActionBookNo,
		}

		butt1, err := json.Marshal(b1)
		if err != nil {
			return nil, err
		}
		butt2, err := json.Marshal(b2)
		if err != nil {
			return nil, err
		}

		button1 := tgbotapi.NewInlineKeyboardButtonData("Занять", string(butt1))
		button2 := tgbotapi.NewInlineKeyboardButtonData("К списку мест", string(butt2))
		row := tgbotapi.NewInlineKeyboardRow(button1, button2)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
		msg.ReplyMarkup = keyboard
	}
	return &msg, nil
}

func (s *messageFormerImpl) FormBookSeatResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {

	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = result.Message

	return &msg, nil
}

func (s *messageFormerImpl) FormCancelBookResult(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = result.Message

	return &msg, nil
}

func (s *messageFormerImpl) FormTextMsg(ctx context.Context, result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	chatID := model.GetCurrentChatID(ctx)

	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = result.Message

	return &msg, nil
}

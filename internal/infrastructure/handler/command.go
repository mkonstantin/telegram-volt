package handler

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
	"telegram-api/internal/infrastructure/handler/dto"
)

type CommandHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type commandHandlerImpl struct {
	userService usecase.UserService
	logger      *zap.Logger
}

func NewCommandHandler(userService usecase.UserService, logger *zap.Logger) CommandHandler {
	return &commandHandlerImpl{
		userService: userService,
		logger:      logger,
	}
}

func (s *commandHandlerImpl) Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "start":
		return s.handleStartCommand(update)
	default:
		msg.Text = "I don't know that command"
	}

	return &msg, nil
}

func (s *commandHandlerImpl) handleStartCommand(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	data := usecasedto.FirstStartDTO{
		User: model.User{
			Name:         update.Message.From.FirstName,
			TelegramID:   update.Message.From.ID,
			TelegramName: update.Message.From.UserName,
		},
		MessageID: update.Message.MessageID,
		ChatID:    update.Message.Chat.ID,
	}

	result, err := s.userService.FirstCome(data)
	if err != nil || result == nil {
		return nil, err
	}

	switch result.Key {
	case usecase.OfficeMenu:
		return sendOfficeMenu(result)
	case usecase.ChooseOfficeMenu:
		return sendChooseOfficeMenu(result)
	}

	// TODO
	return nil, nil
}

func sendOfficeMenu(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
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
	row := tgbotapi.NewInlineKeyboardRow(button1, button2, button3)
	confirmOfficeKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg.Text = result.Message
	msg.ReplyMarkup = confirmOfficeKeyboard
	msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

func sendChooseOfficeMenu(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
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
	msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

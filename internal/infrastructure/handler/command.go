package handler

import (
	"encoding/json"
	"fmt"
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
	case usecase.ChooseOffice:
		return chooseOffice(result)
	case usecase.ConfirmOffice:
		return confirmAlreadyChosenOffice(result)
	}

	// TODO
	return nil, nil
}

func confirmAlreadyChosenOffice(result *usecasedto.FirstStartResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")

	trueAnswer := &dto.CommandResponse{
		Type:      usecase.ConfirmOffice,
		OfficeID:  result.Office.ID,
		IsConfirm: true,
	}
	falseAnswer := &dto.CommandResponse{
		Type:      usecase.ConfirmOffice,
		OfficeID:  result.Office.ID,
		IsConfirm: false,
	}

	trueA, err := json.Marshal(trueAnswer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	falseA, err := json.Marshal(falseAnswer)
	if err != nil {
		return nil, err
	}

	buttonPos := tgbotapi.NewInlineKeyboardButtonData("Да", string(trueA))
	buttonNeg := tgbotapi.NewInlineKeyboardButtonData("Нет", string(falseA))
	row := tgbotapi.NewInlineKeyboardRow(buttonPos, buttonNeg)
	confirmOfficeKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	msg.Text = result.Message
	msg.ReplyMarkup = confirmOfficeKeyboard
	msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

func chooseOffice(result *usecasedto.FirstStartResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, office := range result.Offices {
		resp := &dto.CommandResponse{
			Type:     usecase.ChooseOffice,
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

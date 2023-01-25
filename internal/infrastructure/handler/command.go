package handler

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"telegram-api/internal/app/usecase"
	"telegram-api/internal/domain/model"
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
	data := usecase.UserLogicData{
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
		return s.chooseOffice(result)
	case usecase.ConfirmOffice:
		return s.confirmAlreadyChosenOffice(result)
	}

	return nil, nil
}

type Sdf struct {
	Name string
	Ags  int
}

func (s *commandHandlerImpl) confirmAlreadyChosenOffice(result *usecase.UserLogicResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")
	emp := &Sdf{Name: "Kostya", Ags: 37}
	estr, err := json.Marshal(emp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(estr))
	fff := string(estr)
	buttonPos := tgbotapi.NewInlineKeyboardButtonData("Да", fff)
	buttonNeg := tgbotapi.NewInlineKeyboardButtonData("Нет", fff)
	row := tgbotapi.NewInlineKeyboardRow(buttonPos, buttonNeg)
	confirmOfficeKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)

	message := fmt.Sprintf("%s, хотите занять место в: %s?", result.User.Name, result.Office.Name)
	msg.Text = message
	msg.ReplyMarkup = confirmOfficeKeyboard
	msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

func (s *commandHandlerImpl) chooseOffice(result *usecase.UserLogicResult) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(result.ChatID, "")
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, office := range result.Offices {
		button := tgbotapi.NewInlineKeyboardButtonData(office.Name, strconv.FormatInt(office.ID, 10))
		row := tgbotapi.NewInlineKeyboardRow(button)
		rows = append(rows, row)
	}

	var chooseOfficeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)

	message := fmt.Sprintf("Привет, %s! Давай выберем офис)", result.User.Name)
	msg.Text = message
	msg.ReplyMarkup = chooseOfficeKeyboard
	msg.ReplyToMessageID = result.MessageID
	return &msg, nil
}

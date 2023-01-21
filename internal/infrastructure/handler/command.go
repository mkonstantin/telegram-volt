package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"telegram-api/internal/app/use_case"
)

type CommandHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type commandHandlerImpl struct {
	userService use_case.UserService
	logger      *zap.Logger
}

func NewCommandHandler(userService use_case.UserService, logger *zap.Logger) CommandHandler {
	return &commandHandlerImpl{
		userService: userService,
		logger:      logger,
	}
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func (s *commandHandlerImpl) Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "start":
		name := update.Message.From.FirstName
		strName := fmt.Sprintf("Привет, %s! Для начала давай выберем офис)", name)
		msg.Text = strName
		//msg.ReplyMarkup = firstOfficeChoose
		msg.ReplyMarkup = numericKeyboard
		msg.ReplyToMessageID = update.Message.MessageID
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "I don't know that command"
	}

	return &msg, nil
}

package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"telegram-api/internal/app/usecase"
	usecasedto "telegram-api/internal/app/usecase/dto"
	"telegram-api/internal/domain/model"
)

type CommandHandler interface {
	Handle(update tgbotapi.Update) (*tgbotapi.MessageConfig, error)
}

type commandHandlerImpl struct {
	msgFormer   MessageFormer
	userService usecase.UserService
	logger      *zap.Logger
}

func NewCommandHandler(msgFormer MessageFormer, userService usecase.UserService, logger *zap.Logger) CommandHandler {
	return &commandHandlerImpl{
		msgFormer:   msgFormer,
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
		return s.sendOfficeMenu(result)
	case usecase.ChooseOfficeMenu:
		return s.sendChooseOfficeMenu(result)
	}

	// TODO
	return nil, nil
}

func (s *commandHandlerImpl) sendOfficeMenu(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	return s.msgFormer.FormOfficeMenuMsg(result)
}

func (s *commandHandlerImpl) sendChooseOfficeMenu(result *usecasedto.UserResult) (*tgbotapi.MessageConfig, error) {
	return s.msgFormer.FormChooseOfficeMenuMsg(result)
}
